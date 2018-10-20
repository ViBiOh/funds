package model

import (
	"database/sql"

	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
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

// ListAlertsOpened retrieves current Alerts (only one mail sent)
func (f *App) ListAlertsOpened() (alerts []*Alert, err error) {
	rows, err := f.dbConnexion.Query(listAlertsOpenedQuery)
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
			err = errors.WithStack(err)
			return
		}

		alerts = append(alerts, &Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}

// SaveAlert saves Alert
func (f *App) SaveAlert(alert *Alert, tx *sql.Tx) (err error) {
	if alert == nil {
		return errors.New(`cannot save nil`)
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(f.dbConnexion, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(saveAlertQuery, alert.Isin, alert.Score, alert.AlertType); err != nil {
		err = errors.WithStack(err)
	}

	return
}
