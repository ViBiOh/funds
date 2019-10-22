package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	httputils "github.com/ViBiOh/httputils/v2/pkg"
	"github.com/ViBiOh/httputils/v2/pkg/alcotest"
	"github.com/ViBiOh/httputils/v2/pkg/cors"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/opentracing"
	"github.com/ViBiOh/httputils/v2/pkg/owasp"
	"github.com/ViBiOh/httputils/v2/pkg/prometheus"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	opentracingConfig := opentracing.Flags(fs, "tracing")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	fundsConfig := model.Flags(fs, "")
	dbConfig := db.Flags(fs, "db")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	owaspApp := owasp.New(owaspConfig)
	corsApp := cors.New(corsConfig)

	fundApp, err := model.New(fundsConfig, dbConfig)
	if err != nil {
		logger.Fatal(err)
	}

	handler := httputils.ChainMiddlewares(fundApp.Handler(), prometheusApp, opentracingApp, owaspApp, corsApp)

	go fundApp.Start()

	httputils.New(serverConfig).ListenAndServe(handler, httputils.HealthHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})), nil)
}
