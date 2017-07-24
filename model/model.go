package model

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ViBiOh/funds/cache"
	"github.com/ViBiOh/funds/crawler"
	"github.com/ViBiOh/funds/jsonHttp"
)

const refreshDelayInHours = 6

var listRequest = regexp.MustCompile(`^/list$`)
var performanceURL string

type results struct {
	Results interface{} `json:"results"`
}

var cacheRequests = make(chan cache.Request, crawler.MaxConcurrentFetcher)

// Init start cache server routine and init it from crawling
func Init(url string, dbHost string, dbPort int, dbUser string, dbPass string, dbName string) {
	performanceURL = url

	InitCache()
	if dbHost != `` {
		InitDB(dbHost, dbPort, dbUser, dbPass, dbName)
	}
}

// InitCache load cache
func InitCache() {
	go cache.Server(cacheRequests, len(performanceIds))
	go func() {
		refreshCache()
		c := time.Tick(refreshDelayInHours * time.Hour)
		for range c {
			refreshCache()
		}
	}()
}

func refreshCache() {
	log.Print(`Cache refresh - start`)
	defer log.Print(`Cache refresh - end`)

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

	performancesCache := make([]cache.Content, 0)
	for performance := range results {
		performancesCache = append(performancesCache, performance.(cache.Content))
	}

	cache.Push(cacheRequests, performancesCache)
	if db != nil {
		performances := make([]Performance, 0)
		for _, performance := range performancesCache {
			performances = append(performances, performance.(Performance))
		}

		if err := SaveAll(performances, nil); err != nil {
			log.Printf(`Error while saving Performances: %v`, err)
		}
	}
}

// ListPerformances return content of performance cache
func ListPerformances() []Performance {
	performances := make([]Performance, 0, len(performanceIds))
	for perf := range cache.List(cacheRequests) {
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
