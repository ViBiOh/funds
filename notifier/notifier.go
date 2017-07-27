package notifier

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/model"
)

const from = `funds@vibioh.fr`
const name = `Funds App`
const subject = `[Funds] Score level notification`
const notificationInterval = 24 * time.Hour

func getTimer(hour int, minute int) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
	if !nextTime.After(time.Now()) {
		nextTime = nextTime.Add(notificationInterval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func getCurrentAlerts() (map[string]model.Alert, error) {
	if !db.Ping() {
		return make(map[string]model.Alert, 0), nil
	}

	alerts, err := model.AlertsOpened()
	if err != nil {
		return nil, err
	}

	currentAlerts := make(map[string]model.Alert)
	for _, alert := range alerts {
		if _, ok := currentAlerts[alert.Isin]; !ok {
			currentAlerts[alert.Isin] = alert
		}
	}

	return currentAlerts, nil
}

func getPerformancesAbove(score float64, currentAlerts map[string]model.Alert) ([]model.Performance, error) {
	if !db.Ping() {
		return make([]model.Performance, 0), nil
	}

	performances, err := model.PerformanceWithScoreAbove(score)
	if err != nil {
		return nil, err
	}

	performancesToAlert := make([]model.Performance, 0)
	for _, performance := range performances {
		if alert, ok := currentAlerts[performance.Isin]; ok {
			if alert.AlertType != `above` {
				performancesToAlert = append(performancesToAlert, performance)
			}
		} else {
			performancesToAlert = append(performancesToAlert, performance)
		}
	}

	return performancesToAlert, nil
}

func getPerformancesBelow(currentAlerts map[string]model.Alert) ([]model.Performance, error) {
	if !db.Ping() {
		return make([]model.Performance, 0), nil
	}

	performances := make([]model.Performance, 0)

	for _, alert := range currentAlerts {
		if performance, err := model.PerformanceByIsin(alert.Isin); err != nil {
			return nil, err
		} else if performance.Score < alert.Score {
			performances = append(performances, *performance)
		}
	}

	return performances, nil
}

func saveTypedAlerts(score float64, performances []model.Performance, alertType string) error {
	for _, performance := range performances {
		if err := model.SaveAlert(model.Alert{Isin: performance.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return err
		}
	}

	return nil
}

func saveAlerts(score float64, above []model.Performance, below []model.Performance) error {
	if err := saveTypedAlerts(score, above, `above`); err != nil {
		return fmt.Errorf(`Error while saving above alerts: %v`, err)
	}

	if err := saveTypedAlerts(score, below, `below`); err != nil {
		return fmt.Errorf(`Error while saving below alerts: %v`, err)
	}

	return nil
}

func notify(recipients string, score float64) error {
	currentAlerts, err := getCurrentAlerts()
	if err != nil {
		return fmt.Errorf(`Error while getting current alerts: %v`, err)
	}

	above, err := getPerformancesAbove(score, currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting above performances: %v`, err)
	}

	below, err := getPerformancesBelow(currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting below performances: %v`, err)
	}

	if (len(above) > 0 || len(below) > 0) && recipients != `` {
		htmlContent, err := getHTMLContent(score, above, below)

		if err != nil {
			return err
		}

		if mailjet.Ping() {
			if err := mailjet.SendMail(from, name, subject, strings.Split(recipients, `,`), string(htmlContent)); err != nil {
				return err
			}
			log.Printf(`Sending mail notification for %d funds to %s`, len(above)+len(below), recipients)
		}

		if db.Ping() {
			if err := saveAlerts(score, above, below); err != nil {
				return err
			}
		}
	}

	return nil
}

// Start the notifier
func Start(recipients string, score float64, hour int, minute int) {
	timer := getTimer(hour, minute)

	for {
		select {
		case <-timer.C:
			if err := notify(recipients, score); err != nil {
				log.Print(err)
			}
			timer.Reset(notificationInterval)
		}
	}
}
