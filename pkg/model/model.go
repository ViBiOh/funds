package model

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/cron"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httperror"
	"github.com/ViBiOh/httputils/v3/pkg/httpjson"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

const (
	maxConcurrentFetcher = 8
	listPrefix           = "/list"
	alertsPrefix         = "/alerts"
)

// Config of package
type Config struct {
	infos *string
}

// App of package
type App interface {
	Health() bool
	Start()
	Handler() http.Handler
	ListFunds([]Alert) []Fund
	GetFundsAbove(float64, map[string]Alert) ([]*Fund, error)
	GetFundsBelow(map[string]Alert) ([]*Fund, error)
	GetIsinAlert() ([]Alert, error)
	GetCurrentAlerts() (map[string]Alert, error)
	SaveAlert(*Alert, *sql.Tx) error
}

type app struct {
	dbConnexion *sql.DB
	fundsURL    string
	fundsMap    sync.Map
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: flags.New(prefix, "funds").Name("Infos").Default("").Label("Informations URL").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, dbConfig db.Config) (App, error) {
	app := &app{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},
	}

	fundsDB, err := db.New(dbConfig)
	if err != nil {
		logger.Error("%s", err)
	} else {
		app.dbConnexion = fundsDB
	}

	return app, nil
}

func (a *app) Start() {
	cron.New().Each(time.Hour*8).Now().Start(a.refresh, func(err error) {
		logger.Error("%s", err)
	})
}

func (a *app) refresh(_ time.Time) error {
	if a.fundsURL == "" {
		return nil
	}

	a.refreshData(context.Background())

	if a.dbConnexion != nil {
		if err := a.saveData(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) refreshData(ctx context.Context) {
	inputs := make(chan []byte, maxConcurrentFetcher)

	for i := uint(0); i < maxConcurrentFetcher; i++ {
		go func() {
			for input := range inputs {
				if output, err := fetchFund(ctx, a.fundsURL, input); err != nil {
					logger.Error("%s", err)
				} else {
					a.fundsMap.Store(output.ID, output)
				}
			}
		}()
	}

	for _, fundID := range fundsIds {
		inputs <- fundID
	}
	close(inputs)
}

func (a *app) saveData() (err error) {
	var tx *sql.Tx
	if tx, err = db.GetTx(a.dbConnexion, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err = a.saveFund(&fund, tx)

		return err == nil
	})

	return
}

// Health check health
func (a *app) Health() bool {
	return db.Ping(a.dbConnexion)
}

// ListFunds return content of funds' map
func (a *app) ListFunds(alerts []Alert) []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		for _, alert := range alerts {
			if fund.Isin == alert.Isin {
				fund.Alert = &alert
			}
		}

		funds = append(funds, fund)
		return true
	})

	return funds
}

func (a *app) listHandler(w http.ResponseWriter, r *http.Request) {
	alerts, err := a.GetIsinAlert()
	if err != nil {
		httperror.InternalServerError(w, fmt.Errorf("unable to retrieve alerts: %w", err))
		return
	}

	httpjson.ResponseArrayJSON(w, http.StatusOK, a.ListFunds(alerts), httpjson.IsPretty(r))
}

// Handler for model request. Should be use with net/http
func (a *app) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			if _, err := w.Write(nil); err != nil {
				httperror.InternalServerError(w, err)
			}
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		if strings.HasPrefix(r.URL.Path, listPrefix) {
			a.listHandler(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
