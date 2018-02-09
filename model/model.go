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

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/tools"
)

const (
	maxConcurrentFetcher = 24
	refreshDelay         = 8 * time.Hour
	listPrefix           = `/list`
)

// FundApp wrap all fund methods
type FundApp struct {
	dbConnexion *sql.DB
	fundsURL    string
	fundsMap    sync.Map
}

// NewApp creates FundApp from Flags
func NewApp(config map[string]*string, dbConfig map[string]*string) (*FundApp, error) {
	app := &FundApp{fundsURL: *config[`infos`], fundsMap: sync.Map{}}

	fundsDB, err := db.GetDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf(`Error while initializing database: %v`, err)
	}

	app.dbConnexion = fundsDB

	if app.fundsURL != `` {
		go app.refreshCron()
	}

	return app, nil
}

func (f *FundApp) refreshCron() {
	f.refresh()
	c := time.Tick(refreshDelay)
	for range c {
		f.refresh()
	}
}

func (f *FundApp) refresh() {
	log.Print(`Refresh started`)
	defer log.Print(`Refresh ended`)

	if err := f.refreshData(); err != nil {
		log.Printf(`Error while refreshing: %v`, err)
	}

	if err := f.saveData(); err != nil {
		log.Printf(`Error while saving: %v`, err)
	}
}

func (f *FundApp) refreshData() error {
	inputs, results, errors := tools.ConcurrentAction(maxConcurrentFetcher, func(ID interface{}) (interface{}, error) {
		return fetchFund(f.fundsURL, ID.([]byte))
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
			f.fundsMap.Store(content.ID, content)
			break
		}
	}

	if len(errorIds) > 0 {
		return fmt.Errorf(`Errors with ids %s`, bytes.Join(errorIds, []byte(`,`)))
	}

	return nil
}

func (f *FundApp) saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(f.dbConnexion, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	f.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err = f.SaveFund(&fund, tx)

		return err == nil
	})

	return
}

// Health check health
func (f *FundApp) Health() bool {
	return db.Ping(f.dbConnexion)
}

// ListFunds return content of funds' map
func (f *FundApp) ListFunds() []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	f.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		funds = append(funds, value.(Fund))
		return true
	})

	return funds
}

func (f *FundApp) listHandler(w http.ResponseWriter, r *http.Request) {
	if err := httputils.ResponseArrayJSON(w, http.StatusOK, f.ListFunds(), httputils.IsPretty(r.URL.RawQuery)); err != nil {
		httputils.InternalServerError(w, err)
	}
}

// Flags add flags for given prefix
func Flags(prefix string) map[string]*string {
	return map[string]*string{
		`infos`: flag.String(tools.ToCamel(prefix+`Infos`), ``, `[funds] Informations URL`),
	}
}

// Handler for model request. Should be use with net/http
func Handler(app *FundApp) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			if _, err := w.Write(nil); err != nil {
				httputils.InternalServerError(w, err)
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
