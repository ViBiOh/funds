package main

import (
	"flag"
	"log"
	"os"

	"github.com/ViBiOh/funds/pkg/mailjet"
	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/pkg/db"
)

func main() {
	check := flag.Bool(`c`, false, `Healthcheck (check and exit)`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications in Europe/Paris`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)

	mailjetConfig := mailjet.Flags(``)

	flag.Parse()

	mailjetApp := mailjet.NewApp(mailjetConfig)
	fundApp, err := model.NewApp(fundsConfig, dbConfig)
	if err != nil {
		log.Printf(`Error while creating Fund app: %v`, err)
	}

	if err := notifier.Init(mailjetApp); err != nil {
		log.Printf(`Error while initializing notifier: %v`, err)
	}

	if *check {
		if !fundApp.Health() || !mailjetApp.Ping() {
			os.Exit(1)
		}
		return
	}

	log.Printf(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	notifier.Start(*recipients, *score, *hour, *minute, fundApp)
}
