package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ViBiOh/funds/pkg/model"
	httputils "github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
	"github.com/ViBiOh/httputils/pkg/server"
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

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal("%+v", err)
	}

	alcotest.DoAndExit(alcotestConfig)

	serverApp, err := httputils.New(serverConfig)
	if err != nil {
		logger.Fatal("%+v", err)
	}

	healthcheckApp := healthcheck.New()
	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	gzipApp := gzip.New()
	owaspApp := owasp.New(owaspConfig)
	corsApp := cors.New(corsConfig)

	fundApp, err := model.New(fundsConfig, dbConfig)
	if err != nil {
		logger.Error("%+v", err)
	}

	modelHandler := model.Handler(fundApp)
	healthcheckApp.NextHealthcheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))

	handler := server.ChainMiddlewares(modelHandler, prometheusApp, opentracingApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(handler, nil, healthcheckApp)
}
