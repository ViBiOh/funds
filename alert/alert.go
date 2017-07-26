package main

import (
	"flag"

	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/notifier"
)

func main() {
	dbHost := flag.String(`dbHost`, ``, `Host of Postgres database, leave empty for no database use`)
	dbPort := flag.Int(`dbPort`, 5432, `Port of Postgres database`)
	dbUser := flag.String(`dbUser`, `postgres`, `User of Postgres database`)
	dbPass := flag.String(`dbPass`, `postgres`, `Password of Postgres database`)
	dbName := flag.String(`dbName`, `funds`, `Name of Postgres database`)
	recipients := flag.String(`recipients`, ``, `Email of notifications recipients`)
	score := flag.Float64(`score`, 15.0, `Score value to notification when above`)
	hour := flag.Int(`hour`, 8, `Hour of day for sending notifications`)
	minute := flag.Int(`minute`, 0, `Minute of hour for sending notifications`)
	flag.Parse()

	if *dbHost != `` {
		db.InitDB(*dbHost, *dbPort, *dbUser, *dbPass, *dbName)
	}

	notifier.InitMailjet()

	notifier.Start(*recipients, *score, *hour, *minute)
}
