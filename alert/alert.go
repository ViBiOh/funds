package main

import (
	"flag"

	"github.com/ViBiOh/funds/notifier"
)

func main() {
	apiURL := flag.String(`api`, `https://funds-api.vibioh.fr/list`, `URL of funds-api`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	hourOfDay := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minuteOfHour := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	scoreStep := flag.Float64(`score`, 15.0, `Score value to notification when above`)
	flag.Parse()

	notifier.Start(*apiURL, *recipients, *hourOfDay, *minuteOfHour, *scoreStep)
}
