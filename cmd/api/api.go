package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/cors"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	httputils_model "github.com/ViBiOh/httputils/v3/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	fundsDb, err := db.New(dbConfig)
	logger.Fatal(err)

	fundApp := model.New(fundsConfig, fundsDb)

	go fundApp.Start()

	httputils.New(serverConfig).ListenAndServe(fundApp.Handler(), []httputils_model.Pinger{fundsDb.Ping}, prometheus.New(prometheusConfig).Middleware, owasp.New(owaspConfig).Middleware, cors.New(corsConfig).Middleware)
}
