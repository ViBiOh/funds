package notifier

import (
	"log"
	"time"

	"github.com/ViBiOh/funds/model"
)

const notificationInterval = 24 * time.Hour

func getTimer(hour int, minute int) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
	if !nextTime.After(time.Now()) {
		nextTime = nextTime.Add(notificationInterval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func getPerformances(score float64) ([]model.Performance, error) {
	return model.PerformanceWithScoreAbove(score)
}

func notify(recipients string, score float64) {
	performances, err := getPerformances(score)
	if err != nil {
		log.Printf(`Error while getting performances: %v`, err)
		return
	}

	if len(performances) > 0 {
		htmlContent, err := getHTMLContent(score, performances)

		if err != nil {
			log.Printf(`Error while creating HTML email: %v`, err)
		}

		log.Printf(`%s`, htmlContent)
		log.Printf(`Sended to %s`, recipients)
	}
}

// Start the notifier
func Start(recipients string, score float64, hour int, minute int) {
	timer := getTimer(hour, minute)

	for {
		select {
		case <-timer.C:
			notify(recipients, score)
			timer.Reset(notificationInterval)
		}
	}
}
