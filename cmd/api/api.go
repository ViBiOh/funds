package main

import (
	"flag"
	"net/http"
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

	fundApp, err := model.New(fundsConfig, dbConfig)
	if err != nil {
		logger.Fatal(err)
	}

	go fundApp.Start()

	server := httputils.New(serverConfig)
	server.Health(httputils.HealthHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})))
	server.Middleware(prometheus.New(prometheusConfig))
	server.Middleware(owasp.New(owaspConfig))
	server.Middleware(cors.New(corsConfig))
	server.ListenServeWait(fundApp.Handler())
}
