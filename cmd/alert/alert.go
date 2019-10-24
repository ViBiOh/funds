package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/v2/pkg/cron"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/opentracing"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("alert", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	opentracingConfig := opentracing.Flags(fs, "tracing")
	mailerConfig := client.Flags(fs, "mailer")
	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	if *check {
		return
	}

	opentracing.New(opentracingConfig)

	fundApp, err := model.New(fundsConfig, dbConfig)
	logger.Fatal(err)

	mailerApp := client.New(mailerConfig)

	notifierApp := notifier.New(notifierConfig, fundApp, mailerApp)
	logger.Fatal(err)

	cron.New().Days().At("08:00").In("Europe/Paris").Start(notifierApp.Do, func(err error) {
		logger.Error("%+v", err)
	})
}
