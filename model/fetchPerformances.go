package model

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/ViBiOh/funds/fetch"
)

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

func cleanID(performanceID []byte) string {
	return string(bytes.ToLower(performanceID))
}

func extractLabel(extract *regexp.Regexp, body []byte, defaultValue []byte) []byte {
	if extract == nil {
		return defaultValue
	}

	match := extract.FindSubmatch(body)
	if match == nil || len(match) < 2 {
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

func getPerformance(url string, perf *Performance) error {
	body, err := fetch.GetBody(url + `&tab=1`)
	if err != nil {
		return err
	}

	perf.Isin = string(extractLabel(isinRegex, body, emptyByte))
	perf.Label = string(extractLabel(labelRegex, body, emptyByte))
	perf.Category = string(extractLabel(categoryRegex, body, emptyByte))
	perf.Rating = string(extractLabel(ratingRegex, body, zeroByte))
	perf.OneMonth = extractPerformance(perfOneMonthRegex, body)
	perf.ThreeMonths = extractPerformance(perfThreeMonthRegex, body)
	perf.SixMonths = extractPerformance(perfSixMonthRegex, body)
	perf.OneYear = extractPerformance(perfOneYearRegex, body)

	return nil
}

func getVolatilite(url string, perf *Performance) error {
	body, err := fetch.GetBody(url + `&tab=2`)
	if err != nil {
		return err
	}

	perf.VolThreeYears = extractPerformance(volThreeYearRegex, body)
	return nil
}

func fetchPerformance(performanceID []byte) (Performance, error) {
	cleanID := cleanID(performanceID)
	perf := &Performance{ID: cleanID, Update: time.Now()}

	if err := getPerformance(performanceURL+cleanID, perf); err != nil {
		log.Printf(`Error while fetching performance for %s: %v`, performanceID, err)
		return *perf, err
	}

	if err := getVolatilite(performanceURL+cleanID, perf); err != nil {
		log.Printf(`Error while fetching volatilite for %s: %v`, performanceID, err)
		return *perf, err
	}

	perf.ComputeScore()

	return *perf, nil
}
