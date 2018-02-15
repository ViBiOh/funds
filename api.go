package main

import (
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/funds/model"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/owasp"
)

const healthPrefix = `/health`

func healthHandler(w http.ResponseWriter, r *http.Request, fundApp *model.App) {
	if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func main() {
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)

	httputils.StartMainServer(func() http.Handler {
		fundApp, err := model.NewApp(fundsConfig, dbConfig)
		if err != nil {
			log.Fatalf(`Error while creating Fund app: %v`, err)
		}

		modelHandler := model.Handler(fundApp)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == healthPrefix {
				healthHandler(w, r, fundApp)
			} else {
				modelHandler.ServeHTTP(w, r)
			}
		})

		return gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler)))
	})
}
