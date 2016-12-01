package morningStar

import (
	"../jsonHttp"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const PERFORMANCE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const VOLATILITE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`
const EMPTY_BYTE = []byte(``)
const REFRESH_DELAY = 18

var LIST_REQUEST = regexp.MustCompile(`^/list$`)
var PERF_REQUEST = regexp.MustCompile(`^/(.+?)$`)

var ISIN = regexp.MustCompile(`ISIN.:(\S+)`)
var LABEL = regexp.MustCompile(`<h1[^>]*?>((?:.|\n)*?)</h1>`)
var RATING = regexp.MustCompile(`<span\sclass=".*?stars([0-9]).*?">`)
var CATEGORY = regexp.MustCompile(`<span[^>]*?>Catégorie</span>.*?<span[^>]*?>(.*?)</span>`)
var PERF_ONE_MONTH = regexp.MustCompile(`<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_THREE_MONTH = regexp.MustCompile(`<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_SIX_MONTH = regexp.MustCompile(`<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_ONE_YEAR = regexp.MustCompile(`<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>`)
var VOL_3_YEAR = regexp.MustCompile(`<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>`)

type SyncedMap struct {
	sync.RWMutex
	performances map[string]Performance
}

func (m *SyncedMap) get(key string) (Performance, bool) {
	m.RLock()
	defer m.RUnlock()

	performance, ok := m.performances[key]
	return performance, ok
}

func (m *SyncedMap) push(key string, performance Performance) {
	m.Lock()
	defer m.Unlock()
	m.performances[key] = performance
}

var PERFORMANCE_CACHE = SyncedMap{performances: make(map[string]Performance)}

type Performance struct {
	Id            string    `json:"id"`
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

type Search struct {
	Id    string `json:"i"`
	Label string `json:"n"`
}

type Results struct {
	Results interface{} `json:"results"`
}

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func getBody(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New(`Error while retrieving data from ` + url)
	}

	if response.StatusCode >= 400 {
		return nil, errors.New(`Got error ` + strconv.Itoa(response.StatusCode) + ` while getting ` + url)
	}

	body, err := readBody(response.Body)
	if err != nil {
		return nil, errors.New(`Error while reading body of ` + url)
	}

	return body, nil
}

func getLabel(extract *regexp.Regexp, body []byte, defaultValue []byte) []byte {
	match := extract.FindSubmatch(body)
	if match == nil {
		return defaultValue
	}

	return bytes.Replace(match[1], []byte(`&amp;`), []byte(`&`), -1)
}

func getPerformance(extract *regexp.Regexp, body []byte) float64 {
	dotResult := bytes.Replace(getLabel(extract, body, EMPTY_BYTE), []byte(`,`), []byte(`.`), -1)
	percentageResult := bytes.Replace(dotResult, []byte(`%`), EMPTY_BYTE, -1)
	trimResult := bytes.TrimSpace(percentageResult)

	result, err := strconv.ParseFloat(string(trimResult), 64)
	if err != nil {
		return 0.0
	}
	return result
}

func SinglePerformance(morningStarId []byte) (*Performance, error) {
	cleanId := string(bytes.ToLower(morningStarId))

	performance, ok := PERFORMANCE_CACHE.get(cleanId)

	if ok && time.Now().Add(time.Hour*-REFRESH_DELAY).Before(performance.Update) {
		return &performance, nil
	}

	performanceBody, err := getBody(PERFORMANCE_URL + cleanId)
	if err != nil {
		return nil, err
	}

	volatiliteBody, err := getBody(VOLATILITE_URL + cleanId)
	if err != nil {
		return nil, err
	}

	isin := string(getLabel(ISIN, performanceBody, EMPTY_BYTE))
	label := string(getLabel(LABEL, performanceBody, EMPTY_BYTE))
	rating := string(getLabel(RATING, performanceBody, []byte(`0`)))
	category := string(getLabel(CATEGORY, performanceBody, EMPTY_BYTE))
	oneMonth := getPerformance(PERF_ONE_MONTH, performanceBody)
	threeMonths := getPerformance(PERF_THREE_MONTH, performanceBody)
	sixMonths := getPerformance(PERF_SIX_MONTH, performanceBody)
	oneYear := getPerformance(PERF_ONE_YEAR, performanceBody)
	volThreeYears := getPerformance(VOL_3_YEAR, volatiliteBody)

	score := (0.25 * oneMonth) + (0.3 * threeMonths) + (0.25 * sixMonths) + (0.2 * oneYear) - (0.1 * volThreeYears)
	scoreTruncated := float64(int(score*100)) / 100

	performance = Performance{cleanId, isin, label, category, rating, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears, scoreTruncated, time.Now()}

	PERFORMANCE_CACHE.push(cleanId, performance)

	return &performance, nil
}

func singlePerformanceHandler(w http.ResponseWriter, morningStarId []byte) {
	performance, err := SinglePerformance(morningStarId)

	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		jsonHttp.ResponseJson(w, *performance)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	listBody, err := readBody(r.Body)
	if err != nil {
		http.Error(w, `Error while reading body for list`, 500)
		return
	}

	if len(bytes.TrimSpace(listBody)) == 0 {
		jsonHttp.ResponseJson(w, Results{[0]Performance{}})
		return
	}

	ids := bytes.Split(listBody, []byte(`,`))
	size := len(ids)

	var wg sync.WaitGroup
	wg.Add(size)

	performances := make(chan Performance, size)
	for _, id := range ids {
		go func(morningStarId []byte, ch chan<- Performance) {
			performance, err := SinglePerformance(morningStarId)
			if err == nil {
				ch <- *performance
			}
			wg.Done()
		}(id)
	}

	go func() {
		wg.Wait()
		close(performances)
	}()

	results := make([]Performance, 0, size)
	for performance := range performances {
		results = append(results, performance)
	}

	jsonHttp.ResponseJson(w, Results{results})
}

type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(`Access-Control-Allow-Origin`, `*`)
	w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
	w.Header().Add(`Access-Control-Allow-Methods`, `GET, POST`)
	w.Header().Add(`X-Content-Type-Options`, `nosniff`)

	urlPath := []byte(r.URL.Path)

	if LIST_REQUEST.Match(urlPath) {
		listHandler(w, r)
	} else if PERF_REQUEST.Match(urlPath) {
		singlePerformanceHandler(w, PERF_REQUEST.FindSubmatch(urlPath)[1])
	}
}
