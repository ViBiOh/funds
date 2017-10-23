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

var modelHandler = model.Handler()
var apiHandler http.Handler

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if len(model.ListFunds()) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == `/health` {
			healthHandler(w, r)
		} else {
			modelHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	url := flag.String(`c`, ``, `URL to healthcheck (check and exit)`)
	infosURL := flag.String(`infos`, ``, `Informations URL`)
	dbConfig := db.Flags(``)
	corsConfig := cors.Flags(``)
	flag.Parse()

	if *url != `` {
		alcotest.Do(url)
		return
	}

	fundsDB, err := db.GetDB(dbConfig)
	if err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else if db.Ping(fundsDB) {
		log.Print(`Database ready`)
	}

	if err := model.Init(*infosURL, fundsDB); err != nil {
		log.Printf(`Error while initializing model: %v`, err)
	}

	log.Print(`Starting server on port ` + port)

	apiHandler = prometheus.Handler(`http`, rate.Handler(gziphandler.GzipHandler(owasp.Handler(cors.Handler(corsConfig, handler())))))
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
