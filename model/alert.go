package model

import (
	"fmt"
)

// Alert for a funds
type Alert struct {
	Isin      string
	AlertType string
	Score     float64
}

// GetCurrentAlerts retrieves current opened alerts
func (f *FundApp) GetCurrentAlerts() (map[string]*Alert, error) {
	currentAlerts := make(map[string]*Alert)

	alerts, err := f.ListAlertsOpened()
	if err != nil {
		return nil, fmt.Errorf(`Error while listing opened alerts: %v`, err)
	}

	for _, alert := range alerts {
		if _, ok := currentAlerts[alert.Isin]; !ok {
			currentAlerts[alert.Isin] = alert
		}
	}

	return currentAlerts, nil
}
