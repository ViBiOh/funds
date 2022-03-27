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

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/concurrent"
	"github.com/ViBiOh/httputils/v4/pkg/cron"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/httperror"
	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

const (
	listPrefix = "/list"
)

// Config of package
type Config struct {
	infos *string
}

// App of package
type App struct {
	tracer   trace.Tracer
	db       db.App
	fundsURL string
	fundsMap sync.Map
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: flags.String(fs, prefix, "funds", "Infos", "Informations URL", "", nil),
	}
}

// New creates new App from Config
func New(config Config, db db.App, tracerApp tracer.App) *App {
	return &App{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},

		db:     db,
		tracer: tracerApp.GetTracer("model"),
	}
}

// Start worker
func (a *App) Start(done <-chan struct{}) {
	cron.New().Each(time.Hour*8).Now().WithTracer(a.tracer).OnError(func(err error) {
		logger.Error("%s", err)
	}).Start(a.refresh, done)
}

func (a *App) refresh(ctx context.Context) error {
	if a.fundsURL == "" {
		return nil
	}

	if a.tracer != nil {
		var span trace.Span
		ctx, span = a.tracer.Start(ctx, "refresh")
		defer span.End()
	}

	a.refreshData(ctx)

	if a.db.Enabled() {
		if err := a.saveData(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) refreshData(ctx context.Context) {
	wg := concurrent.NewLimited(4)

	for _, fundID := range fundsIds {
		func(fundID []byte) {
			wg.Go(func() {
				if output, err := fetchFund(ctx, a.fundsURL, fundID); err != nil {
					logger.Error("%s", err)
				} else {
					a.fundsMap.Store(output.ID, output)
				}

				time.Sleep(10 * time.Second)
			})
		}(fundID)
	}

	wg.Wait()
}

func (a *App) saveData(ctx context.Context) (err error) {
	a.fundsMap.Range(func(_ any, value any) bool {
		fund := value.(Fund)
		err := a.db.DoAtomic(ctx, func(ctx context.Context) error {
			return a.saveFund(ctx, &fund)
		})

		return err == nil
	})

	return
}

// Health check health
func (a *App) Health() error {
	if len(a.ListFunds(nil)) == 0 {
		return errors.New("no funds fetched")
	}

	return nil
}

// ListFunds return content of funds' map
func (a *App) ListFunds(alerts []Alert) []Fund {
	funds := make([]Fund, 0, len(fundsIds))

	a.fundsMap.Range(func(_ any, value any) bool {
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

func (a *App) listHandler(w http.ResponseWriter, r *http.Request) {
	alerts, err := a.GetIsinAlert(r.Context())
	if err != nil {
		httperror.InternalServerError(w, fmt.Errorf("unable to retrieve alerts: %w", err))
		return
	}

	httpjson.WriteArray(w, http.StatusOK, a.ListFunds(alerts))
}

// Handler for model request. Should be use with net/http
func (a *App) Handler() http.Handler {
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
