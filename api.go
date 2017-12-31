package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/funds/model"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
)

const port = `1080`

var (
	modelHandler http.Handler
	apiHandler   http.Handler
)

func healthHandler(w http.ResponseWriter, r *http.Request, fundApp *model.FundApp) {
	if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func handler(fundApp *model.FundApp) http.Handler {
	modelHandler = model.Handler(fundApp)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == `/health` {
			healthHandler(w, r, fundApp)
		} else {
			modelHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	alcotestConfig := alcotest.Flags(``)
	prometheusConfig := prometheus.Flags(`prometheus`)
	rateConfig := rate.Flags(`rate`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)
	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	log.Print(`Starting server on port ` + port)

	fundApp, err := model.NewFundApp(fundsConfig, dbConfig)
	if err != nil {
		log.Printf(`Error while creating Fund app: %v`, err)
	}

	apiHandler = prometheus.Handler(prometheusConfig, rate.Handler(rateConfig, gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler(fundApp))))))
	server := &http.Server{
		Addr:    `:` + port,
		Handler: apiHandler,
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		serveError <- server.ListenAndServe()
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
