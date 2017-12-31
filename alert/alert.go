package main

import (
	"flag"
	"log"
	"os"

	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/model"
	"github.com/ViBiOh/funds/notifier"
	"github.com/ViBiOh/httputils/db"
)

func main() {
	check := flag.Bool(`c`, false, `Healthcheck (check and exit)`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications in Europe/Paris`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)
	flag.Parse()

	fundApp, err := model.NewFundApp(fundsConfig, dbConfig)
	if err != nil {
		log.Printf(`Error while creating Fund app: %v`, err)
	}

	if err := notifier.Init(); err != nil {
		log.Printf(`Error while initializing notifier: %v`, err)
	}

	if *check {
		if !fundApp.Health() || !mailjet.Ping() {
			os.Exit(1)
		}
		return
	}

	log.Printf(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	notifier.Start(*recipients, *score, *hour, *minute, fundApp)
}
