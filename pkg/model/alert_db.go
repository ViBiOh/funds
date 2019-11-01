package model

import (
	"database/sql"
	"errors"

	"github.com/ViBiOh/httputils/v3/pkg/db"
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

func (a *app) listAlertsOpened() (alerts []*Alert, err error) {
	rows, err := a.dbConnexion.Query(listAlertsOpenedQuery)
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
		err = rows.Scan(&isin, &alertType, &score)
		if err != nil {
			return
		}

		alerts = append(alerts, &Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}

// SaveAlert saves Alert
func (a *app) SaveAlert(alert *Alert, tx *sql.Tx) (err error) {
	if alert == nil {
		return errors.New("cannot save nil")
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(a.dbConnexion, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	_, err = usedTx.Exec(saveAlertQuery, alert.Isin, alert.Score, alert.AlertType)
	return
}
