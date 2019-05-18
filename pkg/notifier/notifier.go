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
	"github.com/ViBiOh/httputils/pkg/scheduler"
	"github.com/ViBiOh/httputils/pkg/tools"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	from                 = "funds@vibioh.fr"
	name                 = "Funds App"
	subject              = "[Funds] Score level notification"
	notificationInterval = 24 * time.Hour
)

var _ scheduler.Task = &App{}

type scoreTemplateContent struct {
	Score      float64
	AboveFunds []*model.Fund
	BelowFunds []*model.Fund
}

// Config of package
type Config struct {
	mailerURL  *string
	mailerUser *string
	mailerPass *string
	recipients *string
	score      *float64
}

// App of package
type App struct {
	mailerURL  string
	mailerUser string
	mailerPass string
	recipients []string
	score      float64
	modelApp   *model.App
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		mailerURL:  fs.String(tools.ToCamel(fmt.Sprintf("%sMailerURL", prefix)), "", "Mailer URL"),
		mailerUser: fs.String(tools.ToCamel(fmt.Sprintf("%sMailerUser", prefix)), "", "Mailer User"),
		mailerPass: fs.String(tools.ToCamel(fmt.Sprintf("%sMailerPass", prefix)), "", "Mailer Pass"),
		recipients: fs.String(tools.ToCamel(fmt.Sprintf("%sRecipients", prefix)), "", "Email of notifications recipients"),
		score:      fs.Float64(tools.ToCamel(fmt.Sprintf("%sScore", prefix)), 25.0, "Score value to notification when above"),
	}
}

// New creates new App from Config
func New(config Config, modelApp *model.App) *App {
	logger.Info("Notification to %s for score above %.2f", *config.recipients, *config.score)

	return &App{
		mailerURL:  strings.TrimSpace(*config.mailerURL),
		mailerUser: strings.TrimSpace(*config.mailerUser),
		mailerPass: strings.TrimSpace(*config.mailerPass),
		recipients: strings.Split(*config.recipients, ","),
		score:      *config.score,
		modelApp:   modelApp,
	}
}

func (a App) saveTypedAlerts(score float64, funds []*model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := a.modelApp.SaveAlert(&model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}, nil); err != nil {
			return err
		}
	}

	return nil
}

func (a App) saveAlerts(score float64, above []*model.Fund, below []*model.Fund) error {
	if err := a.saveTypedAlerts(score, above, "above"); err != nil {
		return err
	}

	return a.saveTypedAlerts(score, below, "below")
}

// Do send notification to users
func (a App) Do(ctx context.Context, currentTime time.Time) error {
	span := opentracing.StartSpan("Funds Notify")
	defer span.Finish()

	usedCtx := opentracing.ContextWithSpan(ctx, span)

	currentAlerts, err := a.modelApp.GetCurrentAlerts()
	if err != nil {
		return err
	}

	above, err := a.modelApp.GetFundsAbove(a.score, currentAlerts)
	if err != nil {
		return err
	}

	below, err := a.modelApp.GetFundsBelow(currentAlerts)
	if err != nil {
		return err
	}

	if len(a.recipients) > 0 && (len(above) > 0 || len(below) > 0) {
		_, _, _, err := request.DoJSON(usedCtx, fmt.Sprintf("%s/render/funds?from=%s&sender=%s&to=%s&subject=%s", a.mailerURL, url.QueryEscape(from), url.QueryEscape(name), url.QueryEscape(strings.Join(a.recipients, ",")), url.QueryEscape(subject)), scoreTemplateContent{a.score, above, below}, http.Header{"Authorization": []string{request.GenerateBasicAuth(a.mailerUser, a.mailerPass)}}, http.MethodPost)
		if err != nil {
			return err
		}

		logger.Info("Mail notification sent to %d recipients for %d funds", len(a.recipients), len(above)+len(below))

		if err := a.saveAlerts(a.score, above, below); err != nil {
			return err
		}
	}

	return nil
}
