package model

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/tools"
)

const (
	maxConcurrentFetcher = 24
	refreshDelay         = 8 * time.Hour
	listPrefix           = `/list`
)

// App wrap all fund methods
type App struct {
	dbConnexion *sql.DB
	fundsURL    string
	fundsMap    sync.Map
}

// NewApp creates App from Flags
func NewApp(config map[string]*string, dbConfig map[string]*string) (*App, error) {
	app := &App{fundsURL: *config[`infos`], fundsMap: sync.Map{}}

	fundsDB, err := db.GetDB(dbConfig)
	if err != nil {
		return app, fmt.Errorf(`Error while initializing database: %v`, err)
	}

	app.dbConnexion = fundsDB

	if app.fundsURL != `` {
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
	log.Print(`Refresh started`)
	defer log.Print(`Refresh ended`)

	if err := a.refreshData(); err != nil {
		log.Printf(`Error while refreshing: %v`, err)
	}

	if err := a.saveData(); err != nil {
		log.Printf(`Error while saving: %v`, err)
	}
}

func (a *App) refreshData() error {
	inputs, results, errors := tools.ConcurrentAction(maxConcurrentFetcher, func(ID interface{}) (interface{}, error) {
		return fetchFund(a.fundsURL, ID.([]byte))
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
		case crawlErr := <-errors:
			errorIds = append(errorIds, crawlErr.Input.([]byte))
			break
		case result := <-results:
			content := result.(Fund)
			a.fundsMap.Store(content.ID, content)
			break
		}
	}

	if len(errorIds) > 0 {
		return fmt.Errorf(`Errors with ids %s`, bytes.Join(errorIds, []byte(`,`)))
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
	if err := httpjson.ResponseArrayJSON(w, http.StatusOK, a.ListFunds(), httpjson.IsPretty(r.URL.RawQuery)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

// Flags add flags for given prefix
func Flags(prefix string) map[string]*string {
	return map[string]*string{
		`infos`: flag.String(tools.ToCamel(fmt.Sprintf(`%sInfos`, prefix)), ``, `[funds] Informations URL`),
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
