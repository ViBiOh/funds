package model

import (
	"context"
	"errors"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const listLastAlertByIsin = `
SELECT
  isin,
  type,
  score,
  creation_date
FROM
  alerts
WHERE
  (isin, creation_date) IN (
    SELECT
      isin,
      max(creation_date)
    FROM
      alerts
    GROUP BY
      isin
   )
`

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

func (a *app) listLastAlertByIsin() (alerts []Alert, err error) {
	rows, err := a.db.Query(listLastAlertByIsin)
	if err != nil {
		return nil, err
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
		var date time.Time

		err = rows.Scan(&isin, &alertType, &score, &date)
		if err != nil {
			return
		}

		alerts = append(alerts, Alert{Isin: isin, AlertType: alertType, Score: score, Date: date})
	}

	return
}

func (a *app) listAlertsOpened() (alerts []Alert, err error) {
	rows, err := a.db.Query(listAlertsOpenedQuery)
	if err != nil {
		return nil, err
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

		alerts = append(alerts, Alert{Isin: isin, AlertType: alertType, Score: score})
	}

	return
}

// SaveAlert saves Alert
func (a *app) SaveAlert(ctx context.Context, alert *Alert) (err error) {
	if alert == nil {
		return errors.New("cannot save nil")
	}

	return db.Exec(ctx, a.db, saveAlertQuery, alert.Isin, alert.Score, alert.AlertType)
}
