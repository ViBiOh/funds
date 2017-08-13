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

// Init initialize notifier tools
func Init() error {
	if err := InitEmail(); err != nil {
		return err
	}

	return nil
}

func getTimer(hour int, minute int, interval time.Duration) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.UTC)
	if !nextTime.After(time.Now()) {
		nextTime = nextTime.Add(interval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func getCurrentAlerts() (map[string]*model.Alert, error) {
	currentAlerts := make(map[string]*model.Alert)

	if !db.Ping() {
		return currentAlerts, nil
	}

	alerts, err := model.ReadAlertsOpened()
	if err != nil {
		return nil, fmt.Errorf(`Error while reading opened alerts: %v`, err)
	}

	for _, alert := range alerts {
		if _, ok := currentAlerts[alert.Isin]; !ok {
			currentAlerts[alert.Isin] = alert
		}
	}

	return currentAlerts, nil
}

func getFundsAbove(score float64, currentAlerts map[string]*model.Alert) ([]*model.Fund, error) {
	fundsToAlert := make([]*model.Fund, 0)

	if !db.Ping() {
		return fundsToAlert, nil
	}

	funds, err := model.ReadFundsWithScoreAbove(score)
	if err != nil {
		return nil, fmt.Errorf(`Error while reading funds with score >= %.2f: %v`, score, err)
	}

	for _, fund := range funds {
		if alert, ok := currentAlerts[fund.Isin]; ok {
			if alert.AlertType != `above` {
				fundsToAlert = append(fundsToAlert, fund)
			}
		} else {
			fundsToAlert = append(fundsToAlert, fund)
		}
	}

	return fundsToAlert, nil
}

func getFundsBelow(currentAlerts map[string]*model.Alert) ([]*model.Fund, error) {
	funds := make([]*model.Fund, 0)

	if !db.Ping() {
		return funds, nil
	}

	for _, alert := range currentAlerts {
		if fund, err := model.ReadFundByIsin(alert.Isin); err != nil {
			return nil, fmt.Errorf(`Error while reading funds with isin '%s': %v`, alert.Isin, err)
		} else if fund.Score < alert.Score {
			funds = append(funds, fund)
		}
	}

	return funds, nil
}

func saveTypedAlerts(score float64, funds []*model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := model.SaveAlert(&model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return fmt.Errorf(`Error while saving %s alerts: %v`, alertType, err)
		}
	}

	return nil
}

func saveAlerts(score float64, above []*model.Fund, below []*model.Fund) error {
	if err := saveTypedAlerts(score, above, `above`); err != nil {
		return err
	}

	if err := saveTypedAlerts(score, below, `below`); err != nil {
		return err
	}

	return nil
}

func notify(recipients []string, score float64) error {
	currentAlerts, err := getCurrentAlerts()
	if err != nil {
		return fmt.Errorf(`Error while getting current alerts: %v`, err)
	}

	above, err := getFundsAbove(score, currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting above funds: %v`, err)
	}

	below, err := getFundsBelow(currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting below funds: %v`, err)
	}

	if len(recipients) > 0 {
		htmlContent, err := getHTMLContent(score, above, below)
		if err != nil {
			return fmt.Errorf(`Error while generating HTML email: %v`, err)
		}

		if htmlContent == nil {
			return nil
		}

		if mailjet.Ping() {
			if err := mailjet.SendMail(from, name, subject, recipients, string(htmlContent)); err != nil {
				return fmt.Errorf(`Error while sending Mailjet mail: %v`, err)
			}
			log.Printf(`Mail notification sent to %d recipients for %d funds`, len(recipients), len(above)+len(below))
		}

		if db.Ping() {
			if err := saveAlerts(score, above, below); err != nil {
				return fmt.Errorf(`Error while saving alerts: %v`, err)
			}
		}
	}

	return nil
}

// Start the notifier
func Start(recipients string, score float64, hour int, minute int) {
	timer := getTimer(hour, minute, notificationInterval)

	recipientsList := strings.Split(recipients, `,`)

	for {
		select {
		case <-timer.C:
			if err := notify(recipientsList, score); err != nil {
				log.Print(err)
			}
			timer.Reset(notificationInterval)
			log.Printf(`Next notification in %v`, notificationInterval)
		}
	}
}
