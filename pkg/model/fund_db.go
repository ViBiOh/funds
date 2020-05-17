package model

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const fundByIsinQuery = `
SELECT
  label,
  score
FROM
  funds
WHERE
  isin = $1
`

const fundsWithScoreAboveQuery = `
SELECT
  isin,
  label,
  score
FROM
  funds
WHERE
  score >= $1
ORDER BY
  isin ASC
`

const fundsCreateQuery = `
INSERT INTO
  funds
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
  funds
SET
  score = $1,
  update_date = $2
WHERE
  isin = $3
`

var errNilFund = errors.New("unable to save nil Fund")

func scanFunds(rows *sql.Rows, pageSize uint) ([]*Fund, error) {
	var (
		isin  string
		label string
		score float64
	)

	list := make([]*Fund, 0, pageSize)

	for rows.Next() {
		if err := rows.Scan(&isin, &label, &score); err != nil {
			return nil, err
		}

		list = append(list, &Fund{Isin: isin, Label: label, Score: score})
	}

	return list, nil
}

func (a *app) readFundByIsin(ctx context.Context, isin string) (Fund, error) {
	item := Fund{Isin: isin}

	scanner := func(row db.RowScanner) error {
		return row.Scan(&item.Label, &item.Score)
	}
	err := db.GetRow(ctx, a.db, scanner, fundByIsinQuery, isin)

	return item, err
}

func (a *app) listFundsWithScoreAbove(minScore float64) (funds []*Fund, err error) {
	rows, err := a.db.Query(fundsWithScoreAboveQuery, minScore)
	if err != nil {
		return
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanFunds(rows, 0)
}

func (a *app) saveFund(ctx context.Context, fund *Fund) (err error) {
	if fund == nil {
		return errNilFund
	}

	if _, err = a.readFundByIsin(ctx, fund.Isin); err != nil {
		if err == sql.ErrNoRows {
			err = db.Exec(ctx, fundsCreateQuery, fund.Isin, fund.Label, fund.Score)
		}
	} else {
		err = db.Exec(ctx, fundsUpdateScoreQuery, fund.Score, "now()", fund.Isin)
	}

	return
}
