package model

import (
	"github.com/ViBiOh/funds/db"
)

// AlertsByIsin retrieve Alerts by isin
func AlertsByIsin(isin string) ([]Alert, error) {
	rows, err := db.DB.Query(`SELECT type FROM alerts WHERE isin=$1 ORDER BY creation_date DESC`, isin)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var (
		alerts    []Alert
		alertType string
	)

	for rows.Next() {
		if err := rows.Scan(&alertType); err != nil {
			return nil, err
		}

		alerts = append(alerts, Alert{Isin: isin, AlertType: alertType})
	}

	return alerts, nil
}
