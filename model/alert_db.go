package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ViBiOh/httputils/db"
)

const listAlertsOpenedQuery = `
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

const saveAlertQuery = `
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

var errNilAlert = errors.New(`Unable to save nil Alert`)

// ListAlertsOpened retrieves current Alerts (only one mail sent)
func ListAlertsOpened() (alerts []*Alert, err error) {
	rows, err := fundsDB.Query(listAlertsOpenedQuery)
	if err != nil {
		return
	}

	defer func() {
		err = db.RowsClose(rows, err)
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

		alerts = append(alerts, &Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}

// SaveAlert saves Alert
func SaveAlert(alert *Alert, tx *sql.Tx) (err error) {
	if alert == nil {
		return errNilAlert
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(fundsDB, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(saveAlertQuery, alert.Isin, alert.Score, alert.AlertType); err != nil {
		err = fmt.Errorf(`Error while querying: %v`, err)
	}

	return
}
