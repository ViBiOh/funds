package model

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/tools"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	maxConcurrentFetcher = 24
	refreshDelay         = 8 * time.Hour
	listPrefix           = "/list"
)

// Config of package
type Config struct {
	infos *string
}

// App of package
type App struct {
	dbConnexion *sql.DB
	fundsURL    string
	fundsMap    sync.Map
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: fs.String(tools.ToCamel(fmt.Sprintf("%sInfos", prefix)), "", "[funds] Informations URL"),
	}
}

// New creates new App from Config
func New(config Config, dbConfig db.Config) (*App, error) {
	app := &App{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},
	}

	fundsDB, err := db.New(dbConfig)
	if err != nil {
		logger.Error("%+v", errors.WithStack(err))
	} else {
		app.dbConnexion = fundsDB
	}

	if app.fundsURL != "" {
		go app.refreshCron()
	}

	return app, nil
}

func (a *App) refreshCron() {
	a.refresh()
	c := time.Tick(refreshDelay)
	for range c {
		a.refresh()
	}
}

func (a *App) refresh() {
	logger.Info("Refresh started")
	defer logger.Info("Refresh ended")

	if err := a.refreshData(); err != nil {
		logger.Error("%+v", err)
	}

	if a.dbConnexion != nil {
		if err := a.saveData(); err != nil {
			logger.Error("%+v", err)
		}
	}
}

func (a *App) refreshData() error {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Fetch Funds")
	defer span.Finish()

	inputs, results, errs := tools.ConcurrentAction(maxConcurrentFetcher, func(ID interface{}) (interface{}, error) {
		return fetchFund(ctx, a.fundsURL, ID.([]byte))
	})

	go func() {
		defer close(inputs)

		for _, fundID := range fundsIds {
			inputs <- fundID
		}
	}()

	errorIds := make([][]byte, 0)

	for i := 0; i < len(fundsIds); i++ {
		select {
		case crawlErr := <-errs:
			errorIds = append(errorIds, crawlErr.Input.([]byte))
			break
		case result := <-results:
			content := result.(Fund)
			a.fundsMap.Store(content.ID, content)
			break
		}
	}

	if len(errorIds) > 0 {
		return errors.New("errors with ids %s", bytes.Join(errorIds, []byte(",")))
	}

	return nil
}

func (a *App) saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(a.dbConnexion, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err = a.SaveFund(&fund, tx)

		return err == nil
	})

	return
}

// Health check health
func (a *App) Health() bool {
	return db.Ping(a.dbConnexion)
}

// ListFunds return content of funds' map
func (a *App) ListFunds() []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		funds = append(funds, value.(Fund))
		return true
	})

	return funds
}

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	if err := httpjson.ResponseArrayJSON(w, http.StatusOK, a.ListFunds(), httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

// Handler for model request. Should be use with net/http
func Handler(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			if _, err := w.Write(nil); err != nil {
				httperror.InternalServerError(w, err)
			}
			return
		}

		if strings.HasPrefix(r.URL.Path, listPrefix) {
			if r.Method == http.MethodGet {
				app.listHandler(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
