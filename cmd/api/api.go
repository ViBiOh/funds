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
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	fundsDb, err := db.New(dbConfig)
	logger.Fatal(err)

	fundApp := model.New(fundsConfig, fundsDb)

	go fundApp.Start()

	server := httputils.New(serverConfig)
	server.Health(fundsDb.Ping)
	server.Health(fundApp.Health)
	server.Middleware(prometheus.New(prometheusConfig).Middleware)
	server.Middleware(owasp.New(owaspConfig).Middleware)
	server.Middleware(cors.New(corsConfig).Middleware)
	server.ListenServeWait(fundApp.Handler())
}
