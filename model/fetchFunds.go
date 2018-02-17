package model

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ViBiOh/httputils/request"
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

func cleanID(fundID []byte) string {
	return string(bytes.ToLower(fundID))
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

func fetchInfosAndPerformances(url string, fund *Fund) error {
	body, err := request.GetRequest(url+`&tab=1`, nil)
	if err != nil {
		return fmt.Errorf(`Error while fetching: %v`, err)
	}

	fund.Isin = string(extractLabel(isinRegex, body, emptyByte))
	fund.Label = string(extractLabel(labelRegex, body, emptyByte))
	fund.Category = string(extractLabel(categoryRegex, body, emptyByte))
	fund.Rating = string(extractLabel(ratingRegex, body, zeroByte))
	fund.OneMonth = extractPerformance(perfOneMonthRegex, body)
	fund.ThreeMonths = extractPerformance(perfThreeMonthRegex, body)
	fund.SixMonths = extractPerformance(perfSixMonthRegex, body)
	fund.OneYear = extractPerformance(perfOneYearRegex, body)

	return nil
}

func fetchVolatilite(url string, fund *Fund) error {
	body, err := request.GetRequest(url+`&tab=2`, nil)
	if err != nil {
		return fmt.Errorf(`Error while fetching: %v`, err)
	}

	fund.VolThreeYears = extractPerformance(volThreeYearRegex, body)
	return nil
}

func fetchFund(fundsURL string, fundID []byte) (Fund, error) {
	cleanID := cleanID(fundID)
	url := fundsURL + cleanID
	fund := &Fund{ID: cleanID}

	if err := fetchInfosAndPerformances(url, fund); err != nil {
		return *fund, fmt.Errorf(`[%s] Error while fetching infos and performances: %v`, fundID, err)
	}

	if err := fetchVolatilite(url, fund); err != nil {
		return *fund, fmt.Errorf(`[%s] Error while fetching volatilite: %v`, fundID, err)
	}

	fund.ComputeScore()

	return *fund, nil
}
