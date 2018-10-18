package notifier

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/request"
	"github.com/ViBiOh/httputils/pkg/tools"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	from                 = `funds@vibioh.fr`
	name                 = `Funds App`
	subject              = `[Funds] Score level notification`
	notificationInterval = 24 * time.Hour
)

type scoreTemplateContent struct {
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
	locationStr := strings.TrimSpace(*config[`timezone`])
	location, err := time.LoadLocation(locationStr)
	if err != nil {
		return nil, fmt.Errorf(`error while loading location %s: %v`, locationStr, err)
	}

	return &App{
		mailerURL:  strings.TrimSpace(*config[`mailerURL`]),
		mailerUser: strings.TrimSpace(*config[`mailerUser`]),
		mailerPass: strings.TrimSpace(*config[`mailerPass`]),
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

func (a App) getTimer(hour int, minute int, interval time.Duration) *time.Timer {
	nextTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, a.location)
	if !nextTime.After(time.Now().In(a.location)) {
		nextTime = nextTime.Add(interval)
	}

	logger.Info(`Next notification at %v`, nextTime)

	return time.NewTimer(nextTime.Sub(time.Now()))
}

func (a App) saveTypedAlerts(score float64, funds []*model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := a.modelApp.SaveAlert(&model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return fmt.Errorf(`error while saving %s alerts: %v`, alertType, err)
		}
	}

	return nil
}

func (a App) saveAlerts(score float64, above []*model.Fund, below []*model.Fund) error {
	if err := a.saveTypedAlerts(score, above, `above`); err != nil {
		return err
	}

	return a.saveTypedAlerts(score, below, `below`)
}

func (a App) notify(recipients []string, score float64) error {
	span := opentracing.StartSpan(`Funds Notify`)
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	currentAlerts, err := a.modelApp.GetCurrentAlerts()
	if err != nil {
		return fmt.Errorf(`error while getting current alerts: %v`, err)
	}

	above, err := a.modelApp.GetFundsAbove(score, currentAlerts)
	if err != nil {
		return fmt.Errorf(`error while getting above funds: %v`, err)
	}

	below, err := a.modelApp.GetFundsBelow(currentAlerts)
	if err != nil {
		return fmt.Errorf(`error while getting below funds: %v`, err)
	}

	if len(recipients) > 0 && (len(above) > 0 || len(below) > 0) {
		_, err := request.DoJSON(ctx, fmt.Sprintf(`%s/render/funds?from=%s&sender=%s&to=%s&subject=%s`, a.mailerURL, url.QueryEscape(from), url.QueryEscape(name), url.QueryEscape(strings.Join(recipients, `,`)), url.QueryEscape(subject)), scoreTemplateContent{score, above, below}, http.Header{`Authorization`: []string{request.GenerateBasicAuth(a.mailerUser, a.mailerPass)}}, http.MethodPost)
		if err != nil {
			return fmt.Errorf(`error while sending email: %v`, err)
		}

		logger.Info(`Mail notification sent to %d recipients for %d funds`, len(recipients), len(above)+len(below))

		if err := a.saveAlerts(score, above, below); err != nil {
			return fmt.Errorf(`error while saving alerts: %v`, err)
		}
	}

	return nil
}

// Start the notifier
func (a App) Start(recipients string, score float64, hour int, minute int) {
	timer := a.getTimer(hour, minute, notificationInterval)

	recipientsList := strings.Split(recipients, `,`)

	for {
		select {
		case <-timer.C:
			if err := a.notify(recipientsList, score); err != nil {
				logger.Error(`%v`, err)
			}
			timer.Reset(notificationInterval)
			logger.Info(`Next notification in %v`, notificationInterval)
		}
	}
}
