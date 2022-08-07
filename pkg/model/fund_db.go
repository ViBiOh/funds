package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const fundByIsinQuery = `
SELECT
  label,
  score
FROM
  funds.funds
WHERE
  isin = $1
`

const fundsWithScoreAboveQuery = `
SELECT
  isin,
  label,
  score
FROM
  funds.funds
WHERE
  score >= $1
ORDER BY
  isin ASC
`

const fundsCreateQuery = `
INSERT INTO
  funds.funds
(
  isin,
  label,
  score
) VALUES (
  $1,
  $2,
  $3
)
`

const fundsUpdateScoreQuery = `
UPDATE
  funds.funds
SET
  score = $1,
  update_date = $2
WHERE
  isin = $3
`

var errNilFund = errors.New("save nil Fund")

func (a *App) readFundByIsin(ctx context.Context, isin string) (Fund, error) {
	item := Fund{Isin: isin}

	scanner := func(row pgx.Row) error {
		return row.Scan(&item.Label, &item.Score)
	}
	err := a.db.Get(ctx, scanner, fundByIsinQuery, isin)

	return item, err
}

func (a *App) listFundsWithScoreAbove(ctx context.Context, minScore float64) (funds []Fund, err error) {
	list := make([]Fund, 0)
	scanner := func(rows pgx.Rows) error {
		var item Fund

		if err := rows.Scan(&item.Isin, &item.Label, &item.Score); err != nil {
			return fmt.Errorf("scan data: %s", err)
		}

		list = append(list, item)
		return nil
	}

	return list, a.db.List(ctx, scanner, fundsWithScoreAboveQuery, minScore)
}

func (a *App) saveFund(ctx context.Context, fund *Fund) (err error) {
	if fund == nil {
		return errNilFund
	}

	if _, err = a.readFundByIsin(ctx, fund.Isin); err != nil {
		if err == pgx.ErrNoRows {
			err = a.db.Exec(ctx, fundsCreateQuery, fund.Isin, fund.Label, fund.Score)
		}
	} else {
		err = a.db.One(ctx, fundsUpdateScoreQuery, fund.Score, "now()", fund.Isin)
	}

	return
}
