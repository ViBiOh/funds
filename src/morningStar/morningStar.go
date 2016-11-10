package morningStar

import "net/http"
import "errors"
import "time"
import "strings"
import "strconv"
import "regexp"
import "io"
import "io/ioutil"
import "sync"
import "encoding/json"
import "../jsonHttp"

const PERFORMANCE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const VOLATILITE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`
const SEARCH_ID = `http://www.morningstar.fr/fr/util/SecuritySearch.ashx?q=`
const REFRESH_DELAY = 18

var LIST_REQUEST = regexp.MustCompile(`^list$`)
var PERF_REQUEST = regexp.MustCompile(`^(.+?)$`)
var ISIN_REQUEST = regexp.MustCompile(`^(.+?)/isin$`)

var CARRIAGE_RETURN = regexp.MustCompile(`\r?\n`)
var END_CARRIAGE_RETURN = regexp.MustCompile(`\r?\n$`)
var PIPE = regexp.MustCompile(`[|]`)

var ISIN = regexp.MustCompile(`ISIN.:(\S+)`)
var LABEL = regexp.MustCompile(`<h1[^>]*?>((?:.|\n)*?)</h1>`)
var RATING = regexp.MustCompile(`<span\sclass=".*?stars([0-9]).*?">`)
var CATEGORY = regexp.MustCompile(`<span[^>]*?>Catégorie</span>.*?<span[^>]*?>(.*?)</span>`)
var PERF_ONE_MONTH = regexp.MustCompile(`<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_THREE_MONTH = regexp.MustCompile(`<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_SIX_MONTH = regexp.MustCompile(`<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_ONE_YEAR = regexp.MustCompile(`<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>`)
var VOL_3_YEAR = regexp.MustCompile(`<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>`)

var PERFORMANCE_CACHE = struct {
  sync.RWMutex
  m map[string]Performance
}{m: make(map[string]Performance)}

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
  VolThreeYears float64   `json:"v1y"`
  Score         float64   `json:"score"`
  Update        time.Time `json:"ts"`
}

type PerformanceAsync struct {
  performance *Performance
  err         error
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

func getLabel(extract *regexp.Regexp, body []byte) string {
  match := extract.FindSubmatch(body)
  if match == nil {
    return ``
  }

  return string(match[1][:])
}

func getPerformance(extract *regexp.Regexp, body []byte) float64 {
  dotResult := strings.Replace(getLabel(extract, body), `,`, `.`, -1)
  percentageResult := strings.Replace(dotResult, `%`, ``, -1)
  trimResult := strings.TrimSpace(percentageResult)

  result, err := strconv.ParseFloat(trimResult, 64)
  if err != nil {
    return 0.0
  }
  return result
}

func singlePerformance(morningStarId string) (*Performance, error) {
  PERFORMANCE_CACHE.RLock()
  performance, present := PERFORMANCE_CACHE.m[morningStarId]
  PERFORMANCE_CACHE.RUnlock()

  if present && time.Now().Add(time.Hour*-REFRESH_DELAY).Before(performance.Update) {
    return &performance, nil
  }

  performanceBody, err := getBody(PERFORMANCE_URL + morningStarId)
  if err != nil {
    return nil, err
  }

  volatiliteBody, err := getBody(VOLATILITE_URL + morningStarId)
  if err != nil {
    return nil, err
  }

  isin := getLabel(ISIN, performanceBody)
  label := strings.Replace(getLabel(LABEL, performanceBody), `&amp;`, `&`, -1)
  rating := getLabel(RATING, performanceBody)
  category := strings.Replace(getLabel(CATEGORY, performanceBody), `&amp;`, `&`, -1)
  oneMonth := getPerformance(PERF_ONE_MONTH, performanceBody)
  threeMonths := getPerformance(PERF_THREE_MONTH, performanceBody)
  sixMonths := getPerformance(PERF_SIX_MONTH, performanceBody)
  oneYear := getPerformance(PERF_ONE_YEAR, performanceBody)
  volThreeYears := getPerformance(VOL_3_YEAR, volatiliteBody)

  score := (0.25 * oneMonth) + (0.3 * threeMonths) + (0.25 * sixMonths) + (0.2 * oneYear) - (0.1 * volThreeYears)
  scoreTruncated := float64(int(score*100)) / 100

  performance = Performance{morningStarId, isin, label, category, rating, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears, scoreTruncated, time.Now()}

  PERFORMANCE_CACHE.Lock()
  PERFORMANCE_CACHE.m[morningStarId] = performance
  PERFORMANCE_CACHE.Unlock()

  return &performance, nil
}

func singlePerformanceAsync(morningStarId string, ch chan<- PerformanceAsync) {
  performance, err := singlePerformance(morningStarId)
  ch <- PerformanceAsync{performance, err}
}

func singlePerformanceHandler(w http.ResponseWriter, morningStarId string) {
  performance, err := singlePerformance(morningStarId)

  if err != nil {
    http.Error(w, err.Error(), 500)
  } else {
    jsonHttp.ResponseJson(w, *performance)
  }
}

func isinHandler(w http.ResponseWriter, isin string) {
  searchBody, err := getBody(SEARCH_ID + strings.ToLower(isin))
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  cleanBody := END_CARRIAGE_RETURN.ReplaceAllString(string(searchBody[:]), ``)
  lines := CARRIAGE_RETURN.Split(cleanBody, -1)
  size := len(lines)

  results := make([]Search, size)
  for i, line := range lines {
    err := json.Unmarshal([]byte(PIPE.Split(line, -1)[1]), &results[i])
    if err != nil {
      http.Error(w, `Error while unmarshalling data for ISIN `+isin, 500)
    }
  }

  jsonHttp.ResponseJson(w, Results{results})
}

func listHandler(w http.ResponseWriter, r *http.Request) {
  listBody, err := readBody(r.Body)
  if err != nil {
    http.Error(w, `Error while reading body for list`, 500)
    return
  }

  stringBody := string(listBody[:])
  if strings.TrimSpace(stringBody) == `` {
    jsonHttp.ResponseJson(w, Results{[0]Performance{}})
    return
  }

  ids := strings.Split(string(listBody[:]), `,`)
  size := len(ids)

  ch := make(chan PerformanceAsync, size)
  results := make([]Performance, size)

  for _, id := range ids {
    go singlePerformanceAsync(id, ch)
  }

  for i, _ := range ids {
    performanceAsync := <-ch
    if performanceAsync.err == nil {
      results[i] = *performanceAsync.performance
    }
  }

  jsonHttp.ResponseJson(w, Results{results})
}

func Handler(w http.ResponseWriter, r *http.Request) {
  path := strings.ToLower(r.URL.Path)

  w.Header().Add(`Access-Control-Allow-Origin`, `*`)
  w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
  w.Header().Add(`Access-Control-Allow-Methods`, `GET, POST`)
  w.Header().Add(`X-Content-Type-Options`, `nosniff`)

  if LIST_REQUEST.MatchString(path) {
    listHandler(w, r)
  } else if ISIN_REQUEST.MatchString(path) {
    isinHandler(w, ISIN_REQUEST.FindStringSubmatch(path)[1])
  } else if PERF_REQUEST.MatchString(path) {
    singlePerformanceHandler(w, path)
  }
}
