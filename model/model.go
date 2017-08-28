package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/tools"
)

const maxConcurrentFetcher = 32
const refreshDelay = 8 * time.Hour
const listPrefix = `/list`

var fundURL string
var fundsMap *tools.ConcurrentMap

// Init start concurrent map and init it from crawling
func Init(url string) error {
	fundURL = url
	fundsMap = tools.CreateConcurrentMap(len(fundsIds), maxConcurrentFetcher)

	go func() {
		refresh()
		c := time.Tick(refreshDelay)
		for range c {
			refresh()
		}
	}()

	return nil
}

// Shutdown close opened ressource
func Shutdown() error {
	fundsMap.Close()

	return nil
}

func refresh() error {
	log.Printf(`Refresh started`)
	defer log.Printf(`Refresh ended`)

	if err := refreshData(); err != nil {
		log.Printf(`Error while refreshing: %v`, err)
	}

	if db.Ping() {
		if err := saveData(); err != nil {
			log.Printf(`Error while saving: %v`, err)
		}
	}

	return nil
}

func refreshData() error {
	log.Printf(`Data refresh started`)
	defer log.Printf(`Data refresh ended`)

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
			log.Print(crawlErr.Err)
			errorIds = append(errorIds, crawlErr.Input.([]byte))
			break
		case result := <-results:
			content := result.(Fund)
			fundsMap.Push(&content)
			break
		}
	}

	if len(errorIds) > 0 {
		return fmt.Errorf(`Errors with ids %s`, bytes.Join(errorIds, []byte(`,`)))
	}

	return nil
}

const dataSaveLabel = `data save`

func saveData() (err error) {
	log.Printf(`Data save started`)
	defer log.Printf(`Data save ended`)

	var tx *sql.Tx
	if tx, err = db.GetTx(dataSaveLabel, nil); err != nil {
		return err
	}

	defer func() {
		err = db.EndTx(dataSaveLabel, tx, err)
	}()

	for entry := range fundsMap.List() {
		if err == nil {
			err = SaveFund(entry.(*Fund), tx)
		}
	}

	return
}

// ListFunds return content of funds' map
func ListFunds() []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	for entry := range fundsMap.List() {
		funds = append(funds, *(entry.(*Fund)))
	}

	return funds
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseArrayJSON(w, ListFunds())
}

// Handler for model request. Should be use with net/http
type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
}
