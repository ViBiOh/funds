package main

import (
	"flag"
	"log"
	"os"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/notifier"
)

func healthcheck() bool {
	return db.Ping() && mailjet.Ping()
}

func main() {
	check := flag.Bool(`c`, false, `Healthcheck (check and exit)`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 6, `Hour of day for sending notifications in UTC`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	flag.Parse()

	if err := db.Init(); err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else {
		log.Print(`Database ready`)
	}

	if err := mailjet.Init(); err != nil {
		log.Printf(`Error while initializing mailjet: %v`, err)
	} else {
		log.Print(`Mailjet ready`)
	}

	if *check {
		if !healthcheck() {
			os.Exit(1)
		}
		return
	}

	log.Printf(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	if err := notifier.Init(); err != nil {
		log.Printf(`Error while initializing notifier: %v`, err)
	}
	notifier.Start(*recipients, *score, *hour, *minute)
}
