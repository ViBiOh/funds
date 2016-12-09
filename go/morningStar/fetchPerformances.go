package morningStar

import (
	"bytes"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const urlPerformance = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const urlVolatilite = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`
const fetchCount = 2

var emptyByte = []byte(``)
var zeroByte = []byte(`0`)
var periodByte = []byte(`.`)
var commaByte = []byte(`,`)
var percentByte = []byte(`%`)
var ampersandByte = []byte(`&`)
var htmlAmpersandByte = []byte(`&amp;`)

var isinRegex = regexp.MustCompile(`ISIN.:(\S+)`)
var labelRegex = regexp.MustCompile(`\|([^|]*?)\|ISIN`)
var ratingRegex = regexp.MustCompile(`<span\sclass=".*?stars([0-9]).*?">`)
var categoryRegex = regexp.MustCompile(`<span[^>]*?>Catégorie</span>.*?<span[^>]*?>(.*?)</span>`)
var perfOneMonthRegex = regexp.MustCompile(`<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>`)
var perfThreeMonthRegex = regexp.MustCompile(`<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>`)
var perfSixMonthRegex = regexp.MustCompile(`<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>`)
var perfOneYearRegex = regexp.MustCompile(`<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>`)
var volThreeYearRegex = regexp.MustCompile(`<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>`)

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

func getPerformance(wg *sync.WaitGroup, url string, perf *performance, errors chan<- error) {
	defer wg.Done()

	if body, err := getBody(url); err != nil {
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
}

func getVolatilite(wg *sync.WaitGroup, url string, perf *performance, errors chan<- error) {
	defer wg.Done()

	if body, err := getBody(url); err != nil {
		errors <- err
	} else {
		perf.VolThreeYears = extractPerformance(volThreeYearRegex, body)
	}
}

func fetchPerformance(morningStarID []byte) (*performance, error) {
	var wg sync.WaitGroup

	cleanID := cleanID(morningStarID)
	perf := &performance{ID: cleanID, Update: time.Now()}

	wg.Add(fetchCount)
	errors := make(chan error)
	go getPerformance(&wg, urlPerformance+cleanID, perf, errors)
	go getVolatilite(&wg, urlVolatilite+cleanID, perf, errors)

	go func() {
		wg.Wait()
		close(errors)
	}()

	var err error
	for err = range errors {
	}

	perf.computeScore()

	return perf, err
}
