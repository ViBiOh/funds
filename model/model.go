package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ViBiOh/funds/crawler"
	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/jsonHttp"
	"github.com/ViBiOh/funds/tools"
)

const refreshDelay = 8 * time.Hour

var listRequest = regexp.MustCompile(`^/list$`)
var fundURL string
var fundsMap *tools.ConcurrentMap

type results struct {
	Results interface{} `json:"results"`
}

// Init start concurrent map and init it from crawling
func Init(url string) error {
	fundURL = url
	fundsMap = tools.CreateConcurrentMap(len(fundsIds), crawler.MaxConcurrentFetcher)

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

	results, errors := crawler.Crawl(fundsIds, func(ID []byte) (interface{}, error) {
		return fetchFund(ID)
	})

	errorIds := make([][]byte, 0)

	for i := 0; i < len(fundsIds); i++ {
		select {
		case id := <-errors:
			errorIds = append(errorIds, id)
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

func saveData() (err error) {
	log.Printf(`Data save started`)
	defer log.Printf(`Data save ended`)

	var tx *sql.Tx
	if tx, err = db.GetTx(nil); err != nil {
		return err
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	var fund *Fund

	for entry := range fundsMap.List() {
		if err == nil {
			fund = entry.(*Fund)
			err = SaveFund(fund, tx)
		}
	}

	return
}

// ListFunds return content of funds' map
func ListFunds() []Fund {
	var fund *Fund
	funds := make([]Fund, 0, len(fundsIds))

	for entry := range fundsMap.List() {
		fund = entry.(*Fund)
		funds = append(funds, *fund)
	}

	return funds
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	jsonHttp.ResponseJSON(w, results{ListFunds()})
}

// Handler for model request. Should be use with net/http
type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(`Access-Control-Allow-Origin`, `*`)
	w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
	w.Header().Add(`Access-Control-Allow-Methods`, `GET`)
	w.Header().Add(`X-Content-Type-Options`, `nosniff`)

	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	urlPath := []byte(r.URL.Path)

	if listRequest.Match(urlPath) {
		if r.Method == http.MethodGet {
			listHandler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
