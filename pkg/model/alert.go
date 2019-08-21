package model

// Alert for a funds
type Alert struct {
	Isin      string
	AlertType string
	Score     float64
}

// GetCurrentAlerts retrieves current opened alerts
func (a *App) GetCurrentAlerts() (map[string]*Alert, error) {
	currentAlerts := make(map[string]*Alert)

	alerts, err := a.ListAlertsOpened()
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
