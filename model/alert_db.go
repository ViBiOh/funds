package model

import (
	"github.com/ViBiOh/funds/db"
)

const currentAlerts = `
SELECT
  isin,
  type,
  score
FROM
  alerts
WHERE
  isin IN (
    SELECT
      isin
    FROM
      alerts
    GROUP BY
      isin
    HAVING
      MOD(COUNT(type), 2) = 1
  )
ORDER BY
  isin          ASC
  creation_date DESC
`

// AlertsOpened retrieve Alerts not closed (score didn't go below)
func AlertsOpened() (alerts []Alert, err error) {
	rows, err := db.DB.Query(currentAlerts)
	defer func() {
		err = rows.Close()
	}()

	if err != nil {
		return
	}

	var (
		isin      string
		alertType string
		score     float64
	)

	for rows.Next() {
		if err = rows.Scan(&isin, &alertType, &score); err != nil {
			return
		}

		alerts = append(alerts, Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}
