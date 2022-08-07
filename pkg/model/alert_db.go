package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
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

func (a *App) listLastAlertByIsin(ctx context.Context) ([]Alert, error) {
	list := make([]Alert, 0)

	scanner := func(rows pgx.Rows) error {
		var item Alert

		if err := rows.Scan(&item.Isin, &item.AlertType, &item.Score, &item.Date); err != nil {
			return fmt.Errorf("scan data: %s", err)
		}

		list = append(list, item)
		return nil
	}

	return list, a.db.List(ctx, scanner, listLastAlertByIsinQuery)
}

func (a *App) listAlertsOpened(ctx context.Context) ([]Alert, error) {
	list := make([]Alert, 0)

	scanner := func(rows pgx.Rows) error {
		var item Alert

		if err := rows.Scan(&item.Isin, &item.AlertType, &item.Score, &item.Date); err != nil {
			return fmt.Errorf("scan data: %s", err)
		}

		list = append(list, item)
		return nil
	}

	return list, a.db.List(ctx, scanner, listAlertsOpenedQuery)
}

// SaveAlert saves Alert
func (a *App) SaveAlert(ctx context.Context, alert *Alert) (err error) {
	if alert == nil {
		return errors.New("cannot save nil")
	}

	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		return a.db.Exec(ctx, saveAlertQuery, alert.Isin, alert.Score, alert.AlertType)
	})
}
