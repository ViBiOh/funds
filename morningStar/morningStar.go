package morningStar

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/ViBiOh/funds/jsonHttp"
)

const refreshDelayInHours = 6
const maxConcurrentFetcher = 32

var listRequest = regexp.MustCompile(`^/list$`)
var performanceRequest = regexp.MustCompile(`^/(.+?)$`)

// Performance of a funds
type Performance struct {
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

func (perf *Performance) computeScore() {
	score := (0.25 * perf.OneMonth) + (0.3 * perf.ThreeMonths) + (0.25 * perf.SixMonths) + (0.2 * perf.OneYear) - (0.1 * perf.VolThreeYears)
	perf.Score = float64(int(score*100)) / 100
}

type results struct {
	Results interface{} `json:"results"`
}

var cacheRequests = make(chan *cacheRequest, maxConcurrentFetcher)

// Init start cache server routine and init it from crawling
func Init() {
	go cacheServer(cacheRequests, len(morningStarIds))
	go func() {
		refreshCache()
		c := time.Tick(refreshDelayInHours * time.Hour)
		for range c {
			refreshCache()
		}
	}()
}

func concurrentRetrievePerformances(ids [][]byte, wg *sync.WaitGroup, performances chan<- *Performance, errors chan<- []byte) {
	tokens := make(chan int, maxConcurrentFetcher)

	clearSemaphores := func() {
		wg.Done()
		<-tokens
	}

	for _, id := range ids {
		tokens <- 1

		go func(morningStarID []byte) {
			defer clearSemaphores()
			perf, err := fetchPerformance(morningStarID)
			if err == nil {
				performances <- perf
			} else {
				errors <- morningStarID
			}
		}(id)
	}
}

func retrievePerformances(ids [][]byte) ([]*Performance, [][]byte) {
	var wgFetch sync.WaitGroup
	wgFetch.Add(len(ids))

	performancesChan := make(chan *Performance, maxConcurrentFetcher)
	errorsChan := make(chan []byte, maxConcurrentFetcher)

	performances := make([]*Performance, 0, len(ids))
	errors := make([][]byte, 0)

	go concurrentRetrievePerformances(ids, &wgFetch, performancesChan, errorsChan)

	var wgChanDrain sync.WaitGroup

	go func() {
		wgFetch.Wait()
		close(performancesChan)
		close(errorsChan)
	}()

	go func() {
		wgChanDrain.Add(1)

		for err := range errorsChan {
			errors = append(errors, err)
		}

		wgChanDrain.Done()
	}()

	for perf := range performancesChan {
		performances = append(performances, perf)
	}

	wgChanDrain.Wait()

	return performances, errors
}

func refreshCache() {
	log.Print(`Cache refresh - start`)
	defer log.Print(`Cache refresh - end`)

	performances, errors := retrievePerformances(morningStarIds)

	if len(errors) > 0 {
		log.Printf(`Errors while refreshing ids %s`, bytes.Join(errors, []byte(`, `)))
	}

	pushCache(cacheRequests, performances)
}

func retrievePerformance(morningStarID []byte) (*Performance, error) {
	perf := getCache(cacheRequests, cleanID(morningStarID))
	if perf != nil {
		return perf, nil
	}

	perf, err := fetchPerformance(morningStarID)
	if err != nil {
		return nil, err
	}

	pushCache(cacheRequests, []*Performance{perf})
	morningStarIds = append(morningStarIds, morningStarID)

	return perf, nil
}

func performanceHandler(w http.ResponseWriter, morningStarID []byte) {
	perf, err := retrievePerformance(morningStarID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		jsonHttp.ResponseJSON(w, *perf)
	}
}

// ListPerformances return content of performance cache
func ListPerformances() []*Performance {
	performances := make([]*Performance, 0, len(morningStarIds))
	for perf := range listCache(cacheRequests) {
		performances = append(performances, perf)
	}

	return performances
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	jsonHttp.ResponseJSON(w, results{ListPerformances()})
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

	if listRequest.Match(urlPath) && r.Method == http.MethodGet {
		listHandler(w, r)
	} else if performanceRequest.Match(urlPath) && r.Method == http.MethodGet {
		performanceHandler(w, performanceRequest.FindSubmatch(urlPath)[1])
	}
}
