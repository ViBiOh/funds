package notifier

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ViBiOh/funds/fetch"
	"github.com/ViBiOh/funds/model"
)

const notificationInterval = 1 * time.Minute

type apiResult struct {
	Results []model.Performance `json:"results"`
}

func getTimer(hour int, minute int) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
	if !nextTime.After(time.Now()) {
		nextTime = nextTime.Add(notificationInterval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func readFunds(api string) ([]model.Performance, error) {
	data, err := fetch.GetBody(api)
	if err != nil {
		log.Printf(`Error while fetching funds from %s: %v`, api, err)
	}

	result := apiResult{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.Results, nil
}

func getFundsWithAboveScore(scoreStep float64, funds []model.Performance) []model.Performance {
	filteredFunds := make([]model.Performance, 0, len(funds))

	for _, fund := range funds {
		if fund.Score >= scoreStep {
			filteredFunds = append(filteredFunds, fund)
		}
	}

	return filteredFunds
}

func notify(api string, recipients string, score float64) {
	funds, err := readFunds(api)
	if err != nil {
		log.Printf(`Error while reading funds from %s: %v`, api, err)
	}

	scoreFunds := getFundsWithAboveScore(score, funds)
	if len(scoreFunds) > 0 {
		htmlContent, err := getHTMLContent(score, scoreFunds)

		if err != nil {
			log.Printf(`Error while creating HTML email: %v`, err)
		}

		log.Printf(`%s`, htmlContent)
		log.Printf(`Sended to %s`, recipients)
	}
}

// Start the notifier
func Start(api string, recipients string, score float64, hour int, minute int) {
	timer := getTimer(hour, minute)

	for {
		select {
		case <-timer.C:
			notify(api, recipients, score)
			timer.Reset(notificationInterval)
		}
	}
}
