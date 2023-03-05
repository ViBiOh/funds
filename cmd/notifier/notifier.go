package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/request"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("notifier", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")
	tracerConfig := tracer.Flags(fs, "tracer")

	mailerConfig := client.Flags(fs, "mailer")
	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	ctx := context.Background()

	tracerApp, err := tracer.New(ctx, tracerConfig)
	logger.Fatal(err)
	defer tracerApp.Close(ctx)
	request.AddTracerToDefaultClient(tracerApp.GetProvider())

	fundsDb, err := db.New(ctx, dbConfig, tracerApp.GetTracer("database"))
	logger.Fatal(err)
	defer fundsDb.Close()

	mailerApp, err := client.New(mailerConfig, nil, tracerApp.GetTracer("mailer"))
	logger.Fatal(err)
	defer mailerApp.Close()

	fundApp := model.New(fundsConfig, fundsDb, tracerApp.GetTracer("funds"))

	notifierApp := notifier.New(notifierConfig, fundApp, mailerApp)
	logger.Fatal(err)

	ctx, end := tracer.StartSpan(ctx, tracerApp.GetTracer("notifier"), "notifier")
	defer end(nil)

	logger.Fatal(notifierApp.Start(ctx))
}
