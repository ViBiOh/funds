package notifier

import (
	"context"
	"flag"
	"strings"
	"time"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/mailer/pkg/client"
)

const (
	from    = "funds@vibioh.fr"
	name    = "Funds App"
	subject = "[Funds] Score level notification"
)

type scoreTemplateContent struct {
	Score      float64       `json:"score"`
	AboveFunds []*model.Fund `json:"aboveFunds"`
	BelowFunds []*model.Fund `json:"belowFunds"`
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
	recipients []string
	score      float64

	modelApp  model.App
	mailerApp client.App
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		recipients: flags.New(prefix, "notifier").Name("Recipients").Default("").Label("Email of notifications recipients").ToString(fs),
		score:      flags.New(prefix, "notifier").Name("Score").Default(25.0).Label("Score value to notification when above").ToFloat64(fs),
	}
}

// New creates new App from Config
func New(config Config, modelApp model.App, mailerApp client.App) *App {
	logger.Info("Notification to %s for score above %.2f", *config.recipients, *config.score)

	return &App{
		recipients: strings.Split(*config.recipients, ","),
		score:      *config.score,
		modelApp:   modelApp,
		mailerApp:  mailerApp,
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
func (a App) Do(currentTime time.Time) error {
	usedCtx := context.Background()

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
		if err := client.NewEmail(a.mailerApp).From(from).As(name).WithSubject(subject).Data(scoreTemplateContent{a.score, above, below}).To(a.recipients...).Template("funds").Send(usedCtx); err != nil {
			return err
		}

		logger.Info("Mail notification sent to %d recipients for %d funds", len(a.recipients), len(above)+len(below))

		if err := a.saveAlerts(a.score, above, below); err != nil {
			return err
		}
	}

	return nil
}
