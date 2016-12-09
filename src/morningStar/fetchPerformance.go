package morningStar

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const urlPerformance = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const urlVolatilite = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`

var emptyByte = []byte(``)
var zeroByte = []byte(`0`)
var periodByte = []byte(`.`)
var commaByte = []byte(`,`)
var percentByte = []byte(`%`)
var ampersandByte = []byte(`&`)
var htmlAmpersandByte = []byte(`&amp;`)

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

	return bytes.Replace(match[1], htmlAmpersandByte, ampersandByte, -1)
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
	var wg sync.WaitGroup
	wg.Add(2)

	cleanID := cleanID(morningStarID)
	perf := &performance{ID: cleanID, Update: time.Now()}
	errors := make(chan error)

	go func(perf *performance, errors chan<- error) {
		defer wg.Done()
		
		if body, err := getBody(urlPerformance + cleanID); err != nil {
			errors <- err
		} else {
			perf.Isin = string(extractLabel(isinRegex, body, emptyByte))
			perf.Label = string(extractLabel(labelRegex, body, emptyByte))
			perf.Category = string(extractLabel(categoryRegex, body, emptyByte))
			perf.Rating = string(extractLabel(ratingRegex, body, zeroByte))
			perf.OneMonth = extractPerformance(perfOneMonthRegex, body)
			perf.ThreeMonths = extractPerformance(perfThreeMonthRegex, body)
			perf.SixMonths = extractPerformance(perfSixMonthRegex, body)
			perf.OneYear = extractPerformance(perfOneYearRegex, body)
		}
	}(perf, errors)

	go func(perf *performance, errors chan<- error) {
		defer wg.Done()
		
		if body, err := getBody(urlVolatilite + cleanID); err != nil {
			errors <- err
		} else {
			perf.VolThreeYears = extractPerformance(volThreeYearRegex, body)
		}
	}(perf, errors)

	go func() {
		wg.Wait()
		close(errors)
	}()
	
	var err error
	for err = range errors {}
	
	perf.computeScore()

	return perf, err
}
