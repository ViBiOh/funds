package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/opentracing"
	"github.com/ViBiOh/httputils/v2/pkg/scheduler"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("alert", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	opentracingConfig := opentracing.Flags(fs, "tracing")
	schedulerConfig := scheduler.Flags(fs, "scheduler")
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
	schedulerApp, err := scheduler.New(schedulerConfig, notifierApp)
	logger.Fatal(err)

	schedulerApp.Start()
}
