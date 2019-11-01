package model

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ViBiOh/httputils/v3/pkg/request"
)

var emptyByte = []byte("")
var zeroByte = []byte("0")
var periodByte = []byte(".")
var commaByte = []byte(",")
var percentByte = []byte("%")
var ampersandByte = []byte("&")
var htmlAmpersandByte = []byte("&amp;")

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

func fetchInfosAndPerformances(ctx context.Context, url string, fund *Fund) error {
	body, _, _, err := request.Get(ctx, fmt.Sprintf("%s&tab=1", url), nil)
	if err != nil {
		return err
	}

	result, err := request.ReadContent(body)
	if err != nil {
		return err
	}

	fund.Isin = string(extractLabel(isinRegex, result, emptyByte))
	fund.Label = string(extractLabel(labelRegex, result, emptyByte))
	fund.Category = string(extractLabel(categoryRegex, result, emptyByte))
	fund.Rating = string(extractLabel(ratingRegex, result, zeroByte))
	fund.OneMonth = extractPerformance(perfOneMonthRegex, result)
	fund.ThreeMonths = extractPerformance(perfThreeMonthRegex, result)
	fund.SixMonths = extractPerformance(perfSixMonthRegex, result)
	fund.OneYear = extractPerformance(perfOneYearRegex, result)

	return nil
}

func fetchVolatilite(ctx context.Context, url string, fund *Fund) error {
	body, _, _, err := request.Get(ctx, fmt.Sprintf("%s&tab=2", url), nil)
	if err != nil {
		return err
	}

	result, err := request.ReadContent(body)
	if err != nil {
		return err
	}

	fund.VolThreeYears = extractPerformance(volThreeYearRegex, result)
	return nil
}

func fetchFund(ctx context.Context, fundsURL string, fundID []byte) (Fund, error) {
	cleanID := cleanID(fundID)
	url := fmt.Sprintf("%s%s", fundsURL, cleanID)
	fund := &Fund{ID: cleanID}

	if err := fetchInfosAndPerformances(ctx, url, fund); err != nil {
		return *fund, fmt.Errorf("unable to fetch infos for %s: %w", fundID, err)
	}

	if err := fetchVolatilite(ctx, url, fund); err != nil {
		return *fund, fmt.Errorf("unable to fetch volatilite for %s: %w", fundID, err)
	}

	fund.ComputeScore()

	return *fund, nil
}
