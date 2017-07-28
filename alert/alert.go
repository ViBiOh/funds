package main

import (
	"flag"
	"log"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/mailjet"
	"github.com/ViBiOh/funds/notifier"
)

func main() {
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 25.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	flag.Parse()

	db.Init()
	mailjet.Init()

	log.Printf(`Notification to %s at %02d:%02d for score above %.2f`, *recipients, *hour, *minute, *score)

	notifier.Start(*recipients, *score, *hour, *minute)
}
