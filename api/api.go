package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/model"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/gzip"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
)

const port = `1080`

var modelHandler = gzip.Handler{Handler: owasp.Handler{Handler: cors.Handler{Handler: model.Handler{}}}}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if len(model.ListFunds()) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func fundsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == `/health` {
		healthHandler(w, r)
	} else {
		modelHandler.ServeHTTP(w, r)
	}
}

func main() {
	url := flag.String(`c`, ``, `URL to healthcheck (check and exit)`)
	infosURL := flag.String(`infos`, ``, `Informations URL`)
	flag.Parse()

	if *url != `` {
		alcotest.Do(url)
		return
	}

	if err := db.Init(); err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else if db.Ping() {
		log.Print(`Database ready`)
	}

	if err := model.Init(*infosURL); err != nil {
		log.Printf(`Error while initializing model: %v`, err)
	}

	log.Print(`Starting server on port ` + port)

	server := &http.Server{
		Addr:    `:` + port,
		Handler: prometheus.NewPrometheusHandler(`http`, http.HandlerFunc(fundsHandler)),
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		serveError <- server.ListenAndServe()
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
