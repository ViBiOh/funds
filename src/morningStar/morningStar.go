package morningStar

import (
	"../jsonHttp"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const urlIds = `https://elasticsearch.vibioh.fr/funds/morningStarId/_search?size=8000`
const refreshDelayInHours = 12
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
	ThreeMonth    float64   `json:"3m"`
	SixMonth      float64   `json:"6m"`
	OneYear       float64   `json:"1y"`
	VolThreeYears float64   `json:"v3y"`
	Score         float64   `json:"score"`
	Update        time.Time `json:"ts"`
}

type results struct {
	Results interface{} `json:"results"`
}

var cacheRequests = make(chan *cacheRequest, maxConcurrentFetcher)

func init() {
	go performanceCacheServer(cacheRequests)
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
	for _, perf := range retrievePerformances(fetchIds()) {
		cacheRequests <- &cacheRequest{value: perf}
	}
}

func fetchIds() [][]byte {
	idsBody, err := getBody(urlIds)
	if err != nil {
		log.Print(err)
		return nil
	}

	idsMatch := idRegex.FindAllSubmatch(idsBody, -1)

	ids := make([][]byte, 0, len(idsMatch))
	for _, match := range idsMatch {
		ids = append(ids, match[1])
	}

	return ids
}

func retrievePerformance(morningStarID []byte) (*performance, error) {
	cleanID := cleanID(morningStarID)

	request := cacheRequest{key: cleanID, ready: make(chan int)}
	cacheRequests <- &request
	<-request.ready

	perf := request.value
	if perf != nil && time.Now().Add(time.Hour*-(refreshDelayInHours+1)).Before(perf.Update) {
		return perf, nil
	}

	perf, err := fetchPerformance(morningStarID)
	if err != nil {
		return nil, err
	}

	cacheRequests <- &cacheRequest{key: cleanID, value: perf}
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
		http.Error(w, err.Error(), 500)
	} else {
		jsonHttp.ResponseJSON(w, *perf)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	request := cacheRequest{ready: make(chan int)}
	cacheRequests <- &request

	<-request.ready
	jsonHttp.ResponseJSON(w, results{request.list})
}

// Handler for MorningStar request. Should be use with net/http
type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(`Access-Control-Allow-Origin`, `*`)
	w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
	w.Header().Add(`Access-Control-Allow-Methods`, `GET`)
	w.Header().Add(`X-Content-Type-Options`, `nosniff`)

	urlPath := []byte(r.URL.Path)

	if requestList.Match(urlPath) {
		listHandler(w, r)
	} else if requestPerf.Match(urlPath) {
		performanceHandler(w, requestPerf.FindSubmatch(urlPath)[1])
	}
}
