package model

import (
	"database/sql"

	"github.com/ViBiOh/funds/db"
)

const alertsOpenedQuery = `
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
  isin          ASC,
  creation_date DESC
`

// ReadAlertsOpened retrieves current Alerts (only one mail sent)
func ReadAlertsOpened() (alerts []Alert, err error) {
	rows, err := db.DB.Query(alertsOpenedQuery)
	if err != nil {
		return
	}

	defer func() {
		if endErr := rows.Close(); err == nil && endErr != nil {
			err = endErr
		}
	}()

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

// SaveAlert saves Alert
func SaveAlert(alert Alert, tx *sql.Tx) (err error) {
	var usedTx *sql.Tx

	if usedTx, err = db.GetTx(tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	_, err = usedTx.Exec(`INSERT INTO alerts (isin, score, type) VALUES ($1, $2, $3)`, alert.Isin, alert.Score, alert.AlertType)

	return
}
