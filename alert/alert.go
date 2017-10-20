package main

import (
	"flag"
	"log"
	"os"

	dbconfig "github.com/ViBiOh/funds/dbconfig"
	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/notifier"
	"github.com/ViBiOh/httputils/db"
)

func main() {
	check := flag.Bool(`c`, false, `Healthcheck (check and exit)`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications in Europe/Paris`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	flag.Parse()

	fundsDB, err := db.GetDB(*dbconfig.Host, *dbconfig.Port, *dbconfig.User, *dbconfig.Pass, *dbconfig.Name)
	if err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else {
		log.Print(`Database ready`)
	}

	if *check {
		if !(db.Ping(fundsDB) && mailjet.Ping()) {
			os.Exit(1)
		}
		return
	}

	log.Printf(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	if err := notifier.Init(fundsDB); err != nil {
		log.Printf(`Error while initializing notifier: %v`, err)
	}
	notifier.Start(*recipients, *score, *hour, *minute)
}
