package model

import (
	"bytes"
	"database/sql"
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
var performanceURL string
var performanceMap *tools.ConcurrentMap

type results struct {
	Results interface{} `json:"results"`
}

// Init start concurrent map and init it from crawling
func Init(url string) {
	performanceURL = url
	performanceMap = tools.CreateConcurrentMap(len(performanceIds), crawler.MaxConcurrentFetcher)

	go func() {
		refreshData()
		c := time.Tick(refreshDelay)
		for range c {
			refreshData()
		}
	}()
}

func refreshData() {
	results, errors := crawler.Crawl(performanceIds, func(ID []byte) (interface{}, error) {
		return fetchPerformance(ID)
	})

	go func() {
		ids := make([][]byte, 0)

		for id := range errors {
			ids = append(ids, id)
		}

		if len(ids) > 0 {
			log.Printf(`Errors while refreshing ids %s`, bytes.Join(ids, []byte(`, `)))
		}
	}()

	for performance := range results {
		performanceMap.Push(performance.(tools.MapContent))
	}

	if db.Ping() {
		if err := saveData(); err != nil {
			log.Printf(`Error while saving data: %v`, err)
		}
	}
}

func saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(nil); err != nil {
		return err
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	for performance := range performanceMap.List() {
		if err == nil {
			err = SavePerformance(performance.(Performance), tx)
		}
	}

	return
}

// ListPerformances return content of performances' map
func ListPerformances() []Performance {
	performances := make([]Performance, 0, len(performanceIds))
	for perf := range performanceMap.List() {
		performances = append(performances, perf.(Performance))
	}

	return performances
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	jsonHttp.ResponseJSON(w, results{ListPerformances()})
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
