package main

import (
	"flag"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/notifier"
)

func main() {
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 15.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	flag.Parse()

	db.InitDB()
	notifier.InitMailjet()

	notifier.Start(*recipients, *score, *hour, *minute)
}
