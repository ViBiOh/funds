package model

import (
	"context"
	"time"
)

// Alert for a funds
type Alert struct {
	Date      time.Time `json:"date,omitempty"`
	AlertType string    `json:"type"`
	Isin      string    `json:"-"`
	Score     float64   `json:"score"`
}

// GetIsinAlert retrieves last alert occured on by isin
func (a *App) GetIsinAlert(ctx context.Context) ([]Alert, error) {
	return a.listLastAlertByIsin(ctx)
}

// GetCurrentAlerts retrieves current opened alerts
func (a *App) GetCurrentAlerts(ctx context.Context) (map[string]Alert, error) {
	currentAlerts := make(map[string]Alert)

	alerts, err := a.listAlertsOpened(ctx)
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		if _, ok := currentAlerts[alert.Isin]; !ok {
			currentAlerts[alert.Isin] = alert
		}
	}

	return currentAlerts, nil
}
