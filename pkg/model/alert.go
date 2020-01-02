package model

import "time"

// Alert for a funds
type Alert struct {
	Isin      string
	AlertType string
	Score     float64
	Date      time.Time `json:"omitempty"`
}

// GetIsinAlert retrieves last alert occured on by isin
func (a *app) GetIsinAlert() ([]Alert, error) {
	return a.listLastAlertByIsin()
}

// GetCurrentAlerts retrieves current opened alerts
func (a *app) GetCurrentAlerts() (map[string]Alert, error) {
	currentAlerts := make(map[string]Alert)

	alerts, err := a.listAlertsOpened()
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
