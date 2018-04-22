package notifier

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/pkg/request"
	"github.com/ViBiOh/httputils/pkg/tools"
)

const (
	from                 = `funds@vibioh.fr`
	name                 = `Funds App`
	subject              = `[Funds] Score level notification`
	notificationInterval = 24 * time.Hour
)

type ScoreTemplateContent struct {
	Score      float64
	AboveFunds []*model.Fund
	BelowFunds []*model.Fund
}

// App stores informations
type App struct {
	mailerURL  string
	mailerUser string
	mailerPass string
	modelApp   *model.App
	location   *time.Location
}

// NewApp creates new App from Flags' config
func NewApp(config map[string]*string, modelApp *model.App) (*App, error) {
	locationStr := *config[`timezone`]
	location, err := time.LoadLocation(locationStr)
	if err != nil {
		return nil, fmt.Errorf(`Error while loading location %s: %v`, locationStr, err)
	}

	return &App{
		mailerURL:  *config[`mailerURL`],
		mailerUser: *config[`mailerUser`],
		mailerPass: *config[`mailerPass`],
		modelApp:   modelApp,
		location:   location,
	}, nil
}

// Flags adds flags for given prefix
func Flags(prefix string) map[string]*string {
	return map[string]*string{
		`timezone`:   flag.String(tools.ToCamel(fmt.Sprintf(`%sTimezone`, prefix)), `Europe/Paris`, `Timezone`),
		`mailerURL`:  flag.String(tools.ToCamel(fmt.Sprintf(`%sMailerURL`, prefix)), ``, `Mailer URL`),
		`mailerUser`: flag.String(tools.ToCamel(fmt.Sprintf(`%sMailerUser`, prefix)), ``, `Mailer User`),
		`mailerPass`: flag.String(tools.ToCamel(fmt.Sprintf(`%sMailerPass`, prefix)), ``, `Mailer Pass`),
	}
}

func (a *App) getTimer(hour int, minute int, interval time.Duration) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, a.location)
	if !nextTime.After(time.Now().In(a.location)) {
		nextTime = nextTime.Add(interval)
	}

	log.Printf(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func (a *App) saveTypedAlerts(score float64, funds []*model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := a.modelApp.SaveAlert(&model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return fmt.Errorf(`Error while saving %s alerts: %v`, alertType, err)
		}
	}

	return nil
}

func (a *App) saveAlerts(score float64, above []*model.Fund, below []*model.Fund) error {
	if err := a.saveTypedAlerts(score, above, `above`); err != nil {
		return err
	}

	return a.saveTypedAlerts(score, below, `below`)
}

func (a *App) notify(recipients []string, score float64) error {
	currentAlerts, err := a.modelApp.GetCurrentAlerts()
	if err != nil {
		return fmt.Errorf(`Error while getting current alerts: %v`, err)
	}

	above, err := a.modelApp.GetFundsAbove(score, currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting above funds: %v`, err)
	}

	below, err := a.modelApp.GetFundsBelow(currentAlerts)
	if err != nil {
		return fmt.Errorf(`Error while getting below funds: %v`, err)
	}

	if len(recipients) > 0 && (len(above) > 0 || len(below) > 0) {
		_, err := request.DoJSON(fmt.Sprintf(`%s/render/funds/`, a.mailerURL), ScoreTemplateContent{score, above, below}, map[string]string{`Authorization`: request.GetBasicAuth(a.mailerUser, a.mailerPass)}, http.MethodPost)
		if err != nil {
			return fmt.Errorf(`Error while sending email: %v`, err)
		}

		log.Printf(`Mail notification sent to %d recipients for %d funds`, len(recipients), len(above)+len(below))

		if err := a.saveAlerts(score, above, below); err != nil {
			return fmt.Errorf(`Error while saving alerts: %v`, err)
		}
	}

	return nil
}

// Start the notifier
func (a *App) Start(recipients string, score float64, hour int, minute int) {
	timer := a.getTimer(hour, minute, notificationInterval)

	recipientsList := strings.Split(recipients, `,`)

	for {
		select {
		case <-timer.C:
			if err := a.notify(recipientsList, score); err != nil {
				log.Print(err)
			}
			timer.Reset(notificationInterval)
			log.Printf(`Next notification in %v`, notificationInterval)
		}
	}
}
