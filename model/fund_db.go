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

const fundsWithScoreAbove = `
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

// ReadFundByIsin retrieves Fund by isin
func ReadFundByIsin(isin string) (*Fund, error) {
	var (
		label string
		score float64
	)
	err := db.DB.QueryRow(fundByIsinQuery, isin).Scan(&label, &score)

	if err != nil {
		return nil, err
	}

	return &Fund{Isin: isin, Label: label, Score: score}, nil
}

// ReadFundsWithScoreAbove retrieves Fund with score above given level
func ReadFundsWithScoreAbove(minScore float64) (funds []Fund, err error) {
	rows, err := db.DB.Query(fundsWithScoreAbove, minScore)
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
		_, err = tx.Exec(`INSERT INTO funds (isin, label, score) VALUES ($1, $2, $3)`, fund.Isin, fund.Label, fund.Score)
	} else {
		_, err = tx.Exec(`UPDATE funds SET score = $1, update_date = $2 WHERE isin = $3`, fund.Score, `now()`, fund.Isin)
	}

	return
}
