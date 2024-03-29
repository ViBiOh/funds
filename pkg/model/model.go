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

type Config struct {
	infos *string
}

type App struct {
	tracer   trace.Tracer
	db       db.App
	fundsURL string
	fundsMap sync.Map
}

func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		infos: flags.String(fs, prefix, "funds", "Infos", "Informations URL", "", nil),
	}
}

func New(config Config, db db.App, tracer trace.Tracer) *App {
	return &App{
		fundsURL: strings.TrimSpace(*config.infos),
		fundsMap: sync.Map{},

		db:     db,
		tracer: tracer,
	}
}

func (a *App) Start(ctx context.Context) {
	cron.New().Each(time.Hour*8).Now().WithTracer(a.tracer).OnError(func(err error) {
		logger.Error("%s", err)
	}).Start(ctx, a.refresh)
}

func (a *App) refresh(ctx context.Context) error {
	if a.fundsURL == "" {
		return nil
	}

	var err error

	ctx, end := tracer.StartSpan(ctx, a.tracer, "refresh")
	defer end(&err)

	a.refreshData(ctx)

	if a.db.Enabled() {
		if err = a.saveData(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) refreshData(ctx context.Context) {
	wg := concurrent.NewLimited(4)

	for _, fundID := range fundsIds {
		fundID := fundID

		wg.Go(func() {
			if output, err := fetchFund(ctx, a.fundsURL, fundID); err != nil {
				logger.Error("%s", err)
			} else {
				a.fundsMap.Store(output.ID, output)
			}

			time.Sleep(10 * time.Second)
		})
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

func (a *App) Health() error {
	if len(a.ListFunds(nil)) == 0 {
		return errors.New("no funds fetched")
	}

	return nil
}

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
		httperror.InternalServerError(w, fmt.Errorf("retrieve alerts: %w", err))
		return
	}

	httpjson.WriteArray(w, http.StatusOK, a.ListFunds(alerts))
}

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
