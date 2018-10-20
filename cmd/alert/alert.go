package main

import (
	"flag"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
)

func main() {
	check := flag.Bool(`c`, false, `Healthcheck (check and exit)`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)

	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)
	notifierConfig := notifier.Flags(``)
	opentracingConfig := opentracing.Flags(`tracing`)

	flag.Parse()

	if *check {
		return
	}

	opentracing.NewApp(opentracingConfig)

	fundApp, err := model.NewApp(fundsConfig, dbConfig)
	if err != nil {
		logger.Error(`%+v`, err)
	}

	notifierApp, err := notifier.NewApp(notifierConfig, fundApp)
	if err != nil {
		logger.Error(`%+v`, err)
	}

	logger.Info(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	notifierApp.Start(*recipients, *score, *hour, *minute)
}
