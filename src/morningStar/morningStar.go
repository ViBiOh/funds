package morningStar

import (
	"../jsonHttp"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const urlIds = `https://elasticsearch.vibioh.fr/funds/morningStarId/_search?size=8000`
const urlPerformance = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const urlVolatilite = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`
const refreshDelayInHours = 12
const maxConcurrentFetcher = 32

var emptyByte = []byte(``)
var zeroByte = []byte(`0`)
var periodByte = []byte(`.`)
var commaByte = []byte(`,`)
var percentByte = []byte(`%`)
var ampersandByte = []byte(`&`)
var htmAmpersandByte = []byte(`&amp;`)

var requestList = regexp.MustCompile(`^/list$`)
var requestPerf = regexp.MustCompile(`^/(.+?)$`)

var idRegex = regexp.MustCompile(`"_id":"(.*?)"`)
var isinRegex = regexp.MustCompile(`ISIN.:(\S+)`)
var labelRegex = regexp.MustCompile(`\|([^|]*?)\|ISIN`)
var ratingRegex = regexp.MustCompile(`<span\sclass=".*?stars([0-9]).*?">`)
var categoryRegex = regexp.MustCompile(`<span[^>]*?>Catégorie</span>.*?<span[^>]*?>(.*?)</span>`)
var perfOneMonthRegex = regexp.MustCompile(`<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>`)
var perfThreeMonthRegex = regexp.MustCompile(`<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>`)
var perfSixMonthRegex = regexp.MustCompile(`<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>`)
var perfOneYearRegex = regexp.MustCompile(`<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>`)
var volThreeYearRegex = regexp.MustCompile(`<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>`)

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

type cacheRequest struct {
	key   string
	value *performance
	list  []*performance
	ready chan int
}

var cacheRequests = make(chan *cacheRequest)

func cacheMonitor() {
	cache := make(map[string]*performance)

	var ready = func(request *cacheRequest) {
		if request.ready != nil {
			close(request.ready)
		}
	}

	for request := range cacheRequests {
		if request.value != nil {
			cache[request.value.ID] = request.value
		} else if request.key != `` {
			request.value, _ = cache[request.key]
			ready(request)
		} else {
			request.list = make([]*performance, 0, len(cache))
			for _, perf := range cache {
				request.list = append(request.list, perf)
			}

			ready(request)
		}
	}
}

func init() {
	go cacheMonitor()
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

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func getBody(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(`Error while retrieving data from %s: %v`, url, err)
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf(`Got error %d while getting %s`, response.StatusCode, url)
	}

	body, err := readBody(response.Body)
	if err != nil {
		return nil, fmt.Errorf(`Error while reading body of %s: %v`, url, err)
	}

	return body, nil
}

func extractLabel(extract *regexp.Regexp, body []byte, defaultValue []byte) []byte {
	match := extract.FindSubmatch(body)
	if match == nil {
		return defaultValue
	}

	return bytes.Replace(match[1], htmAmpersandByte, ampersandByte, -1)
}

func extractPerformance(extract *regexp.Regexp, body []byte) float64 {
	dotResult := bytes.Replace(extractLabel(extract, body, emptyByte), commaByte, periodByte, -1)
	percentageResult := bytes.Replace(dotResult, percentByte, emptyByte, -1)
	trimResult := bytes.TrimSpace(percentageResult)

	result, err := strconv.ParseFloat(string(trimResult), 64)
	if err != nil {
		return 0.0
	}
	return result
}

func cleanID(morningStarID []byte) string {
	return string(bytes.ToLower(morningStarID))
}

func fetchPerformance(morningStarID []byte) (*performance, error) {
	cleanID := cleanID(morningStarID)
	performanceBody, err := getBody(urlPerformance + cleanID)
	if err != nil {
		return nil, err
	}

	volatiliteBody, err := getBody(urlVolatilite + cleanID)
	if err != nil {
		return nil, err
	}

	isin := string(extractLabel(isinRegex, performanceBody, emptyByte))
	label := string(extractLabel(labelRegex, performanceBody, emptyByte))
	rating := string(extractLabel(ratingRegex, performanceBody, zeroByte))
	category := string(extractLabel(categoryRegex, performanceBody, emptyByte))
	oneMonth := extractPerformance(perfOneMonthRegex, performanceBody)
	threeMonths := extractPerformance(perfThreeMonthRegex, performanceBody)
	sixMonths := extractPerformance(perfSixMonthRegex, performanceBody)
	oneYear := extractPerformance(perfOneYearRegex, performanceBody)
	volThreeYears := extractPerformance(volThreeYearRegex, volatiliteBody)

	score := (0.25 * oneMonth) + (0.3 * threeMonths) + (0.25 * sixMonths) + (0.2 * oneYear) - (0.1 * volThreeYears)
	scoreTruncated := float64(int(score*100)) / 100

	return &performance{cleanID, isin, label, category, rating, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears, scoreTruncated, time.Now()}, nil
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

	request := cacheRequest{key: cleanID}
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
