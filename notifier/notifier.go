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

func saveTypedAlerts(ID int, score float64, performances []model.Performance, alertType string) error {
	for _, performance := range performances {
		if err := model.SaveAlert(model.Alert{ID: ID, Isin: performance.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return err
		}
	}

	return nil
}

func saveAlerts(ID int, score float64, above []model.Performance, below []model.Performance) error {
	if err := saveTypedAlerts(ID, score, above, `above`); err != nil {
		return fmt.Errorf(`Error while saving above alerts: %v`, err)
	}

	if err := saveTypedAlerts(ID, score, below, `below`); err != nil {
		return fmt.Errorf(`Error while saving below alerts: %v`, err)
	}

	return nil
}

func notify(recipients string, score float64) {
	above, err := getPerformancesAbove(score)
	if err != nil {
		log.Printf(`Error while getting above performances: %v`, err)
		return
	}

	below, err := getPerformancesBelow()
	if err != nil {
		log.Printf(`Error while getting below performances: %v`, err)
		return
	}

	if (len(above) > 0 || len(below) > 0) && recipients != `` {
		htmlContent, err := getHTMLContent(score, above, below)

		if err != nil {
			log.Printf(`Error while creating HTML email: %v`, err)
		} else if apiPublicKey != `` {
			if ID, err := MailjetSend(from, name, subject, strings.Split(recipients, `,`), string(htmlContent)); err != nil {
				log.Printf(`Error while sending Mailjet mail: %v`, err)
			} else if err = saveAlerts(ID, score, above, below); err != nil {
				log.Print(err)
			}
		}
	}
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
