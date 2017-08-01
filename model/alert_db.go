package model

import (
	"database/sql"
	"fmt"

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

const alertsCreateQuery = `
INSERT INTO
  alerts
(
  isin,
  score,
  type
) VALUES (
  $1,
  $2,
  $3
)
`

// ReadAlertsOpened retrieves current Alerts (only one mail sent)
func ReadAlertsOpened() (alerts []Alert, err error) {
	rows, err := db.Query(alertsOpenedQuery)
	if err != nil {
		err = fmt.Errorf(`Error while querying alerts opened: %v`, err)
		return
	}

	defer func() {
		err = db.RowsClose(rows, err, `alerts opened`)
	}()

	var (
		isin      string
		alertType string
		score     float64
	)

	for rows.Next() {
		if err = rows.Scan(&isin, &alertType, &score); err != nil {
			err = fmt.Errorf(`Error while scanning alerts opened: %v`, err)
			return
		}

		alerts = append(alerts, Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}

// SaveAlert saves Alert
func SaveAlert(alert *Alert, tx *sql.Tx) (err error) {
	if alert == nil {
		return fmt.Errorf(`Unable to save nil Alert`)
	}

	var usedTx *sql.Tx

	if usedTx, err = db.GetTx(tx); err != nil {
		err = fmt.Errorf(`Error while getting transaction for alert creation: %v`, err)
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(alertsCreateQuery, alert.Isin, alert.Score, alert.AlertType); err != nil {
		err = fmt.Errorf(`Error while creating alert for isin=%s: %v`, alert.Isin, err)
	}

	return
}
