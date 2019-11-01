package model

import (
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

func (a *app) readFundByIsin(isin string) (*Fund, error) {
	var (
		label string
		score float64
	)

	err := a.dbConnexion.QueryRow(fundByIsinQuery, isin).Scan(&label, &score)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &Fund{Isin: isin, Label: label, Score: score}, nil
}

func (a *app) listFundsWithScoreAbove(minScore float64) (funds []*Fund, err error) {
	rows, err := a.dbConnexion.Query(fundsWithScoreAboveQuery, minScore)
	if err != nil {
		return
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanFunds(rows, 0)
}

func (a *app) saveFund(fund *Fund, tx *sql.Tx) (err error) {
	if fund == nil {
		return errNilFund
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

	if _, err = a.readFundByIsin(fund.Isin); err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec(fundsCreateQuery, fund.Isin, fund.Label, fund.Score)
		}
	} else {
		_, err = tx.Exec(fundsUpdateScoreQuery, fund.Score, "now()", fund.Isin)
	}

	return
}
