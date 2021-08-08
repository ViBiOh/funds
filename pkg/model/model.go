package model

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/cron"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/httperror"
	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

const (
	maxConcurrentFetcher = 1
	listPrefix           = "/list"
)

// Config of package
type Config struct {
	infos *string
}

// App of package
type App interface {
	Health() error
	Start(<-chan struct{})
	Handler() http.Handler
	ListFunds([]Alert) []Fund
	GetFundsAbove(float64, map[string]Alert) ([]Fund, error)
	GetFundsBelow(map[string]Alert) ([]Fund, error)
	GetIsinAlert() ([]Alert, error)
	GetCurrentAlerts() (map[string]Alert, error)
	SaveAlert(context.Context, *Alert) error
}

type app struct {
	db       db.App
	fundsURL string
	fundsMap sync.Map
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: flags.New(prefix, "funds").Name("Infos").Default("").Label("Informations URL").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, db db.App) App {
	return &app{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},

		db: db,
	}
}

func (a *app) Start(done <-chan struct{}) {
	cron.New().Each(time.Hour*8).Now().OnError(func(err error) {
		logger.Error("%s", err)
	}).Start(a.refresh, done)
}

func (a *app) refresh(ctx context.Context) error {
	if a.fundsURL == "" {
		return nil
	}

	a.refreshData(ctx)

	if a.db.Enabled() {
		if err := a.saveData(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) refreshData(ctx context.Context) {
	inputs := make(chan []byte, maxConcurrentFetcher)

	go func() {
		for i := uint(0); i < maxConcurrentFetcher; i++ {
			for input := range inputs {
				if output, err := fetchFund(ctx, a.fundsURL, input); err != nil {
					logger.Error("%s", err)
				} else {
					a.fundsMap.Store(output.ID, output)
				}

				time.Sleep(10 * time.Second)
			}
		}
	}()

	for _, fundID := range fundsIds {
		inputs <- fundID
	}
	close(inputs)
}

func (a *app) saveData() (err error) {
	ctx := context.Background()

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		err := a.db.DoAtomic(ctx, func(ctx context.Context) error {
			return a.saveFund(ctx, &fund)
		})

		return err == nil
	})

	return
}

// Health check health
func (a *app) Health() error {
	if len(a.ListFunds(nil)) == 0 {
		return errors.New("no funds fetched")
	}

	return nil
}

// ListFunds return content of funds' map
func (a *app) ListFunds(alerts []Alert) []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	a.fundsMap.Range(func(_ interface{}, value interface{}) bool {
		fund := value.(Fund)
		for _, alert := range alerts {
			fundAlert := alert

			if fund.Isin == alert.Isin {
				fund.Alert = &fundAlert
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

	httpjson.WriteArray(w, http.StatusOK, a.ListFunds(alerts), httpjson.IsPretty(r))
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
		} else if r.URL.Path == "/ready" {
			if a.Health() == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusTeapot)
			}
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
