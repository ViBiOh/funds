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

const maxConcurrentFetcher = 10
const refreshDelay = 8 * time.Hour
const listPrefix = `/list`

var fundURL = flag.String(`infos`, ``, `Informations URL`)
var dbConfig = db.Flags(`db`)
var fundsDB *sql.DB

var fundsMap = sync.Map{}

// Init start concurrent map and init it from crawling
func Init() (err error) {
	fundsDB, err = db.GetDB(dbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing database: %v`, err)
	}

	if *fundURL != `` {
		go func() {
			refresh()
			c := time.Tick(refreshDelay)
			for range c {
				refresh()
			}
		}()
	}

	return
}

// Health check health
func Health() bool {
	return db.Ping(fundsDB)
}

func refresh() error {
	log.Print(`Refresh started`)
	defer log.Print(`Refresh ended`)

	if err := refreshData(); err != nil {
		log.Printf(`Error while refreshing: %v`, err)
	}

	if err := saveData(); err != nil {
		log.Printf(`Error while saving: %v`, err)
	}

	return nil
}

func refreshData() error {
	inputs, results, errors := tools.ConcurrentAction(maxConcurrentFetcher, func(ID interface{}) (interface{}, error) {
		return fetchFund(ID.([]byte))
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
			fundsMap.Store(content.ID, content)
			break
		}
	}

	if len(errorIds) > 0 {
		return fmt.Errorf(`Errors with ids %s`, bytes.Join(errorIds, []byte(`,`)))
	}

	return nil
}

func saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(fundsDB, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err = SaveFund(&fund, tx)

		return err == nil
	})

	return
}

// ListFunds return content of funds' map
func ListFunds() []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	fundsMap.Range(func(_ interface{}, value interface{}) bool {
		funds = append(funds, value.(Fund))
		return true
	})

	return funds
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseArrayJSON(w, http.StatusOK, ListFunds(), httputils.IsPretty(r.URL.RawQuery))
}

// Handler for model request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Write(nil)
			return
		}

		if strings.HasPrefix(r.URL.Path, listPrefix) {
			if r.Method == http.MethodGet {
				listHandler(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
