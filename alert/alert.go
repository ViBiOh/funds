package main

import (
	"flag"

	"github.com/ViBiOh/funds/notifier"
)

func main() {
	api := flag.String(`api`, `https://funds-api.vibioh.fr/list`, `URL of funds-api`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	score := flag.Float64(`score`, 15.0, `Score value to notification when above`)
	flag.Parse()

	notifier.Start(*api, *recipients, *score, *hour, *minute)
}
