package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const listLastAlertByIsinQuery = `
SELECT
  isin,
  type,
  score,
  creation_date
FROM
  funds.alerts
WHERE
  (isin, creation_date) IN (
    SELECT
      isin,
      max(creation_date)
    FROM
      funds.alerts
    GROUP BY
      isin
   )
`

const listAlertsOpenedQuery = `
SELECT
  isin,
  type,
  score,
  creation_date
FROM
  funds.alerts
WHERE
  isin IN (
    SELECT
      isin
    FROM
      funds.alerts
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
  funds.alerts
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

func (a *app) listLastAlertByIsin(ctx context.Context) ([]Alert, error) {
	list := make([]Alert, 0)

	scanner := func(rows *sql.Rows) error {
		var item Alert

		if err := rows.Scan(&item.Isin, &item.AlertType, &item.Score, &item.Date); err != nil {
			return fmt.Errorf("unable to scan data: %s", err)
		}

		list = append(list, item)
		return nil
	}

	return list, db.List(ctx, a.db, scanner, listLastAlertByIsinQuery)
}

func (a *app) listAlertsOpened(ctx context.Context) ([]Alert, error) {
	list := make([]Alert, 0)

	scanner := func(rows *sql.Rows) error {
		var item Alert

		if err := rows.Scan(&item.Isin, &item.AlertType, &item.Score, &item.Date); err != nil {
			return fmt.Errorf("unable to scan data: %s", err)
		}

		list = append(list, item)
		return nil
	}

	return list, db.List(ctx, a.db, scanner, listAlertsOpenedQuery)
}

// SaveAlert saves Alert
func (a *app) SaveAlert(ctx context.Context, alert *Alert) (err error) {
	if alert == nil {
		return errors.New("cannot save nil")
	}

	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		return db.Exec(ctx, saveAlertQuery, alert.Isin, alert.Score, alert.AlertType)
	})
}
