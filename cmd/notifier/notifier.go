package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/funds/pkg/notifier"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("notifier", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	mailerConfig := client.Flags(fs, "mailer")
	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	if *check {
		return
	}

	fundsDb, err := db.New(dbConfig)
	logger.Fatal(err)
	defer func() {
		if err := fundsDb.Close(); err != nil {
			logger.Error("error while closing database connection: %s", err)
		}
	}()

	mailerApp, err := client.New(mailerConfig)
	logger.Fatal(err)
	defer mailerApp.Close()

	fundApp := model.New(fundsConfig, fundsDb)

	notifierApp := notifier.New(notifierConfig, fundApp, mailerApp)
	logger.Fatal(err)

	notifierApp.Start()
}
