package model

import (
	"database/sql"
	"fmt"

	"github.com/ViBiOh/funds/db"
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
)`

const fundsUpdateScoreQuery = `
UPDATE
	funds
SET
	score = $1,
	update_date = $2
WHERE
	isin = $3
`

// ReadFundByIsin retrieves Fund by isin
func ReadFundByIsin(isin string) (*Fund, error) {
	var (
		label string
		score float64
	)
	err := db.QueryRow(fundByIsinQuery, isin).Scan(&label, &score)

	if err != nil {
		return nil, err
	}

	return &Fund{Isin: isin, Label: label, Score: score}, nil
}

// ReadFundsWithScoreAbove retrieves Fund with score above given level
func ReadFundsWithScoreAbove(minScore float64) (funds []Fund, err error) {
	rows, err := db.Query(fundsWithScoreAboveQuery, minScore)
	if err != nil {
		return
	}

	defer func() {
		if endErr := rows.Close(); err == nil && endErr != nil {
			err = endErr
		}
	}()

	var (
		isin  string
		label string
		score float64
	)

	for rows.Next() {
		if err = rows.Scan(&isin, &label, &score); err != nil {
			return
		}

		funds = append(funds, Fund{Isin: isin, Label: label, Score: score})
	}

	return
}

// SaveFund saves Fund
func SaveFund(fund *Fund, tx *sql.Tx) (err error) {
	if fund == nil {
		return fmt.Errorf(`Unable to save nil Fund`)
	}

	var usedTx *sql.Tx

	if usedTx, err = db.GetTx(tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = ReadFundByIsin(fund.Isin); err != nil {
		_, err = tx.Exec(fundsCreateQuery, fund.Isin, fund.Label, fund.Score)
	} else {
		_, err = tx.Exec(fundsUpdateScoreQuery, fund.Score, `now()`, fund.Isin)
	}

	return
}
