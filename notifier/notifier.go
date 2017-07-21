package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/ViBiOh/funds/fetch"
	"github.com/ViBiOh/funds/morningStar"
)

const notificationInterval = 1 * time.Minute

type apiResult struct {
	Results []morningStar.Performance `json:"results"`
}

func getTimer(hour int, minute int) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
	if !nextTime.After(time.Now()) {
		nextTime = nextTime.Add(notificationInterval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func readFunds(apiURL string) ([]morningStar.Performance, error) {
	data, err := fetch.GetBody(apiURL)
	if err != nil {
		log.Printf(`Error while fetching funds from %s: %v`, apiURL, err)
	}

	result := apiResult{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.Results, nil
}

func getFundsWithAboveScore(scoreStep float64, funds []morningStar.Performance) []*morningStar.Performance {
	filteredFunds := make([]*morningStar.Performance, 0, len(funds))

	for _, fund := range funds {
		if fund.Score >= scoreStep {
			filteredFunds = append(filteredFunds, &fund)
		}
	}

	return filteredFunds
}

func main() {
	apiURL := flag.String(`api`, `https://funds-api.vibioh.fr/list`, `URL of funds-api`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	hourOfDay := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minuteOfHour := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	scoreStep := flag.Float64(`score`, 15.0, `Score value to notification when above`)
	flag.Parse()

	timer := getTimer(*hourOfDay, *minuteOfHour)
	for {
		select {
		case <-timer.C:
			funds, err := readFunds(*apiURL)
			if err != nil {
				log.Printf(`Error while reading funds from %s: %v`, *apiURL, err)
			}

			scoreFunds := getFundsWithAboveScore(*scoreStep, funds)
			if len(scoreFunds) > 0 {
				log.Printf(`Sended to %s`, *recipients)
			}
			timer.Reset(notificationInterval)
		}
	}
}
