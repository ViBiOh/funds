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
	"go.opentelemetry.io/otel/trace"
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

	tracerApp, err := tracer.New(tracerConfig)
	logger.Fatal(err)
	defer tracerApp.Close()
	request.AddTracerToDefaultClient(tracerApp.GetProvider())

	fundsDb, err := db.New(dbConfig, tracerApp)
	logger.Fatal(err)
	defer fundsDb.Close()

	mailerApp, err := client.New(mailerConfig, nil)
	logger.Fatal(err)
	defer mailerApp.Close()

	fundApp := model.New(fundsConfig, fundsDb, tracerApp)

	notifierApp := notifier.New(notifierConfig, fundApp, mailerApp)
	logger.Fatal(err)

	ctx := context.Background()
	tracer := tracerApp.GetTracer("notifier")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "notifier")
		defer span.End()
	}

	logger.Fatal(notifierApp.Start(ctx))
}
