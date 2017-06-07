package morningStar

import (
	"github.com/ViBiOh/funds/jsonHttp"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const refreshDelayInHours = 6
const maxConcurrentFetcher = 32

var requestList = regexp.MustCompile(`^/list$`)
var requestPerf = regexp.MustCompile(`^/(.+?)$`)

type performance struct {
	ID            string    `json:"id"`
	Isin          string    `json:"isin"`
	Label         string    `json:"label"`
	Category      string    `json:"category"`
	Rating        string    `json:"rating"`
	OneMonth      float64   `json:"1m"`
	ThreeMonths   float64   `json:"3m"`
	SixMonths     float64   `json:"6m"`
	OneYear       float64   `json:"1y"`
	VolThreeYears float64   `json:"v3y"`
	Score         float64   `json:"score"`
	Update        time.Time `json:"ts"`
}

func (perf *performance) computeScore() {
	score := (0.25 * perf.OneMonth) + (0.3 * perf.ThreeMonths) + (0.25 * perf.SixMonths) + (0.2 * perf.OneYear) - (0.1 * perf.VolThreeYears)
	perf.Score = float64(int(score*100)) / 100
}

type results struct {
	Results interface{} `json:"results"`
}

var cacheRequests = make(chan *cacheRequest, maxConcurrentFetcher)

func init() {
	go cacheServer(cacheRequests)
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

	loadCache(cacheRequests, retrievePerformances(morningStarIds))
}

func retrievePerformance(morningStarID []byte) (*performance, error) {
	perf := getCache(cacheRequests, cleanID(morningStarID))
	if perf != nil {
		return perf, nil
	}

	perf, err := fetchPerformance(morningStarID)
	if err != nil {
		return nil, err
	}

	pushCache(cacheRequests, perf)
	morningStarIds = append(morningStarIds, morningStarID)

	return perf, nil
}

func concurrentRetrievePerformances(ids [][]byte, wg *sync.WaitGroup, performances chan<- *performance) {
	tokens := make(chan int, maxConcurrentFetcher)

	clearSemaphores := func() {
		wg.Done()
		<-tokens
	}

	for _, id := range ids {
		tokens <- 1

		go func(morningStarID []byte) {
			defer clearSemaphores()
			if perf, err := fetchPerformance(morningStarID); err == nil {
				performances <- perf
			}
		}(id)
	}
}

func retrievePerformances(ids [][]byte) []*performance {
	var wg sync.WaitGroup
	wg.Add(len(ids))

	performances := make(chan *performance, maxConcurrentFetcher)
	go concurrentRetrievePerformances(ids, &wg, performances)

	go func() {
		wg.Wait()
		close(performances)
	}()

	results := make([]*performance, 0, len(ids))
	for perf := range performances {
		results = append(results, perf)
	}

	return results
}

func performanceHandler(w http.ResponseWriter, morningStarID []byte) {
	perf, err := retrievePerformance(morningStarID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		jsonHttp.ResponseJSON(w, *perf)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	perfs := make([]*performance, 0, len(morningStarIds))
	for perf := range listCache(cacheRequests) {
		perfs = append(perfs, perf)
	}

	jsonHttp.ResponseJSON(w, results{perfs})
}

// Handler for MorningStar request. Should be use with net/http
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

	if requestList.Match(urlPath) {
		listHandler(w, r)
	} else if requestPerf.Match(urlPath) {
		performanceHandler(w, requestPerf.FindSubmatch(urlPath)[1])
	}
}
