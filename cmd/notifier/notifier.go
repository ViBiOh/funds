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

	mailerConfig := client.Flags(fs, "mailer")
	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")
	notifierConfig := notifier.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	fundsDb, err := db.New(dbConfig)
	logger.Fatal(err)
	defer fundsDb.Close()

	mailerApp, err := client.New(mailerConfig, nil)
	logger.Fatal(err)
	defer mailerApp.Close()

	fundApp := model.New(fundsConfig, fundsDb)

	notifierApp := notifier.New(notifierConfig, fundApp, mailerApp)
	logger.Fatal(err)

	logger.Fatal(notifierApp.Start())
}
