package notifier

import (
	"context"
	"flag"
	"strings"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/mailer/pkg/client"
	mailerModel "github.com/ViBiOh/mailer/pkg/model"
)

const (
	from    = "funds@vibioh.fr"
	name    = "Funds App"
	subject = "[Funds] Score level notification"
)

type scoreTemplateContent struct {
	Score      float64      `json:"score"`
	AboveFunds []model.Fund `json:"aboveFunds"`
	BelowFunds []model.Fund `json:"belowFunds"`
}

// App of package
type App interface {
	Start() error
}

// Config of package
type Config struct {
	recipients *string
	score      *float64
}

// App of package
type app struct {
	modelApp  model.App
	mailerApp client.App

	recipients []string
	score      float64
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		recipients: flags.New(prefix, "notifier").Name("Recipients").Default("").Label("Email of notifications recipients").ToString(fs),
		score:      flags.New(prefix, "notifier").Name("Score").Default(25.0).Label("Score value to notification when above").ToFloat64(fs),
	}
}

// New creates new App from Config
func New(config Config, modelApp model.App, mailerApp client.App) App {
	logger.Info("Notification to %s for score above %.2f", *config.recipients, *config.score)

	return &app{
		recipients: strings.Split(*config.recipients, ","),
		score:      *config.score,
		modelApp:   modelApp,
		mailerApp:  mailerApp,
	}
}

func (a app) saveTypedAlerts(ctx context.Context, score float64, funds []model.Fund, alertType string) error {
	for _, fund := range funds {
		if err := a.modelApp.SaveAlert(ctx, &model.Alert{Isin: fund.Isin, Score: score, AlertType: alertType}); err != nil {
			return err
		}
	}

	return nil
}

func (a app) saveAlerts(ctx context.Context, score float64, above []model.Fund, below []model.Fund) error {
	if err := a.saveTypedAlerts(ctx, score, above, "above"); err != nil {
		return err
	}

	return a.saveTypedAlerts(ctx, score, below, "below")
}

// Start notifier
func (a app) Start() error {
	ctx := context.Background()

	currentAlerts, err := a.modelApp.GetCurrentAlerts()
	if err != nil {
		return err
	}

	above, err := a.modelApp.GetFundsAbove(a.score, currentAlerts)
	if err != nil {
		return err
	}
	logger.Info("Got %d funds above %f", len(above), a.score)

	below, err := a.modelApp.GetFundsBelow(currentAlerts)
	if err != nil {
		return err
	}
	logger.Info("Got %d funds below their initial alert", len(above))

	if len(a.recipients) > 0 && (len(above) > 0 || len(below) > 0) {
		if err := a.mailerApp.Send(ctx, *mailerModel.NewMailRequest().From(from).As(name).WithSubject(subject).Data(scoreTemplateContent{a.score, above, below}).To(a.recipients...).Template("funds")); err != nil {
			return err
		}

		logger.Info("Mail notification sent to %d recipients for %d funds", len(a.recipients), len(above)+len(below))

		if err := a.saveAlerts(ctx, a.score, above, below); err != nil {
			return err
		}
	}

	return nil
}
