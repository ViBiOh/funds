package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
)

func main() {
	fs := flag.NewFlagSet("alert", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")
	recipients := fs.String("recipients", "", "Email of notifications recipients")
	score := fs.Float64("score", 25.0, "Score value to notification when above")
	hour := fs.Int("hour", 8, "Hour of day for sending notifications")
	minute := fs.Int("minute", 0, "Minute of hour for sending notifications")

	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")
	opentracingConfig := opentracing.Flags(fs, "tracing")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal("%+v", err)
	}

	if *check {
		return
	}

	opentracing.New(opentracingConfig)

	fundApp, err := model.New(fundsConfig, dbConfig)
	if err != nil {
		logger.Error("%+v", err)
	}

	notifierApp, err := notifier.New(notifierConfig, fundApp)
	if err != nil {
		logger.Error("%+v", err)
	}

	logger.Info("Notification to %s at %02d:%02d for score above %.2f", *recipients, *hour, *minute, *score)

	notifierApp.Start(*recipients, *score, *hour, *minute)
}
