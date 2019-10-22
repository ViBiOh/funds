package model

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/v2/pkg/cron"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/errors"
	"github.com/ViBiOh/httputils/v2/pkg/httperror"
	"github.com/ViBiOh/httputils/v2/pkg/httpjson"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/tools"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	maxConcurrentFetcher = 8
	listPrefix           = "/list"
)

// Config of package
type Config struct {
	infos *string
}

// App of package
type App interface {
	Health() bool
	Start()
	Handler() http.Handler
	ListFunds() []Fund
	GetFundsAbove(float64, map[string]*Alert) ([]*Fund, error)
	GetFundsBelow(map[string]*Alert) ([]*Fund, error)
	GetCurrentAlerts() (map[string]*Alert, error)
	SaveAlert(*Alert, *sql.Tx) error
	Do(time.Time) error
}

type app struct {
	dbConnexion *sql.DB
	fundsURL    string
	fundsMap    sync.Map
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: tools.NewFlag(prefix, "funds").Name("Infos").Default("").Label("Informations URL").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, dbConfig db.Config) (App, error) {
	app := &app{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},
	}

	fundsDB, err := db.New(dbConfig)
	if err != nil {
		logger.Error("%#v", errors.WithStack(err))
	} else {
		app.dbConnexion = fundsDB
	}

	return app, nil
}

func (a *app) Start() {
	if err := a.Do(time.Now()); err != nil {
		logger.Error("%+v", err)
	}

	cron.NewCron().Each(time.Hour*8).Start(a.Do, func(err error) {
		logger.Error("%+v", err)
	})
}

// Do do scheduler task of refreshing data
func (a *app) Do(_ time.Time) error {
	if a.fundsURL == "" {
		return nil
	}

	logger.Info("Refresh started")
	defer logger.Info("Refresh ended")

	if err := a.refreshData(context.Background()); err != nil {
		logger.Error("%#v", err)
	}

	if a.dbConnexion != nil {
		if err := a.saveData(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) refreshData(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Fetch Funds")
	defer span.Finish()

	inputs, results := tools.ConcurrentAction(maxConcurrentFetcher, true, func(ID interface{}) (interface{}, error) {
		return fetchFund(ctx, a.fundsURL, ID.([]byte))
	})

	go func() {
		defer close(inputs)

		for _, fundID := range fundsIds {
			inputs <- fundID
		}
	}()

	errorIds := make([][]byte, 0)

	for {
		result, ok := <-results

		if !ok {
			break
		}

		if result.Err != nil {
			errorIds = append(errorIds, result.Input.([]byte))
		} else {
			content := result.Output.(Fund)
			a.fundsMap.Store(content.ID, content)
		}
	}

	if len(errorIds) > 0 {
		return errors.New("errors with ids %s", bytes.Join(errorIds, []byte(",")))
	}

	return nil
}

func (a *app) saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(a.dbConnexion, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err = a.saveFund(&fund, tx)

		return err == nil
	})

	return
}

// Health check health
func (a *app) Health() bool {
	return db.Ping(a.dbConnexion)
}

// ListFunds return content of funds' map
func (a *app) ListFunds() []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		funds = append(funds, value.(Fund))
		return true
	})

	return funds
}

func (a *app) listHandler(w http.ResponseWriter, r *http.Request) {
	if err := httpjson.ResponseArrayJSON(w, http.StatusOK, a.ListFunds(), httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

// Handler for model request. Should be use with net/http
func (a *app) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			if _, err := w.Write(nil); err != nil {
				httperror.InternalServerError(w, errors.WithStack(err))
			}
			return
		}

		if strings.HasPrefix(r.URL.Path, listPrefix) {
			if r.Method == http.MethodGet {
				a.listHandler(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
