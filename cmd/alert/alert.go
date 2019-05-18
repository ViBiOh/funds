package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/scheduler"
)

func main() {
	fs := flag.NewFlagSet("alert", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	opentracingConfig := opentracing.Flags(fs, "tracing")
	schedulerConfig := scheduler.Flags(fs, "scheduler")
	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal("%+v", err)
	}

	if *check {
		return
	}

	opentracing.New(opentracingConfig)

	fundApp, err := model.New(fundsConfig, dbConfig)
	if err != nil {
		logger.Fatal("%+v", err)
	}

	notifierApp := notifier.New(notifierConfig, fundApp)
	schedulerApp, err := scheduler.New(schedulerConfig, notifierApp)
	if err != nil {
		logger.Fatal("%+v", err)
	}

	schedulerApp.Start()
}
