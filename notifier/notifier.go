package notifier

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ViBiOh/funds/db"
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

func getPerformancesAbove(score float64) ([]model.Performance, error) {
	if db.DB == nil {
		return make([]model.Performance, 0), nil
	}
	return model.PerformanceWithScoreAbove(score)
}

func getPerformancesBelow() ([]model.Performance, error) {
	if db.DB == nil {
		return make([]model.Performance, 0), nil
	}

	alerts, err := model.AlertsOpened()
	if err != nil {
		return nil, err
	}

	performances := make([]model.Performance, 0)

	for _, alert := range alerts {
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
	above, err := getPerformancesAbove(score)
	if err != nil {
		return fmt.Errorf(`Error while getting above performances: %v`, err)
	}

	below, err := getPerformancesBelow()
	if err != nil {
		return fmt.Errorf(`Error while getting below performances: %v`, err)
	}

	if (len(above) > 0 || len(below) > 0) && recipients != `` {
		htmlContent, err := getHTMLContent(score, above, below)

		if err != nil {
			return fmt.Errorf(`Error while creating HTML email: %v`, err)
		}

		if apiPublicKey != `` {
			if err := MailjetSend(from, name, subject, strings.Split(recipients, `,`), string(htmlContent)); err != nil {
				return fmt.Errorf(`Error while sending Mailjet mail: %v`, err)
			}

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
	notify(recipients, score)

	for {
		select {
		case <-timer.C:
			notify(recipients, score)
			timer.Reset(notificationInterval)
		}
	}
}
