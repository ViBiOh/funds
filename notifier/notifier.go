package notifier

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/model"
)

const (
	locationStr          = `Europe/Paris`
	from                 = `funds@vibioh.fr`
	name                 = `Funds App`
	subject              = `[Funds] Score level notification`
	notificationInterval = 24 * time.Hour
)

var location *time.Location

// Init initialize notifier tools
func Init() (err error) {
	location, err = time.LoadLocation(locationStr)
	if err != nil {
		err = fmt.Errorf(`Error while loading location %s: %v`, locationStr, err)
		return
	}

	if err = InitEmail(); err != nil {
		err = fmt.Errorf(`Error while initializing email: %v`, err)
		return
	}

	return
}

func getTimer(hour int, minute int, interval time.Duration) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, location)
	if !nextTime.After(time.Now().In(location)) {
		nextTime = nextTime.Add(interval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func saveTypedAlerts(fundApp *model.FundApp, score float64, funds []*model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := fundApp.SaveAlert(&model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return fmt.Errorf(`Error while saving %s alerts: %v`, alertType, err)
		}
	}

	return nil
}

func saveAlerts(fundApp *model.FundApp, score float64, above []*model.Fund, below []*model.Fund) error {
	if err := saveTypedAlerts(fundApp, score, above, `above`); err != nil {
		return err
	}

	return saveTypedAlerts(fundApp, score, below, `below`)
}

func notify(fundApp *model.FundApp, recipients []string, score float64) error {
	currentAlerts, err := fundApp.GetCurrentAlerts()
	if err != nil {
		return fmt.Errorf(`Error while getting current alerts: %v`, err)
	}

	above, err := fundApp.GetFundsAbove(score, currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting above funds: %v`, err)
	}

	below, err := fundApp.GetFundsBelow(currentAlerts)
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

		if err := mailjet.SendMail(from, name, subject, recipients, string(htmlContent)); err != nil {
			return fmt.Errorf(`Error while sending Mailjet mail: %v`, err)
		}
		log.Printf(`Mail notification sent to %d recipients for %d funds`, len(recipients), len(above)+len(below))

		if err := saveAlerts(fundApp, score, above, below); err != nil {
			return fmt.Errorf(`Error while saving alerts: %v`, err)
		}
	}

	return nil
}

// Start the notifier
func Start(recipients string, score float64, hour int, minute int, fundApp *model.FundApp) {
	timer := getTimer(hour, minute, notificationInterval)

	recipientsList := strings.Split(recipients, `,`)

	for {
		select {
		case <-timer.C:
			if err := notify(fundApp, recipientsList, score); err != nil {
				log.Print(err)
			}
			timer.Reset(notificationInterval)
			log.Printf(`Next notification in %v`, notificationInterval)
		}
	}
}
