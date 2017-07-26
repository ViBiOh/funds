package model

import (
	"database/sql"

	"github.com/ViBiOh/funds/db"
)

// PerformanceByIsin retrieve Performance by isin
func PerformanceByIsin(isin string) (*Performance, error) {
	var (
		label string
		score float64
	)
	err := db.DB.QueryRow(`SELECT label, score FROM funds WHERE isin = $1`, isin).Scan(&label, &score)

	if err != nil {
		return nil, err
	}

	return &Performance{Isin: isin, Label: label, Score: score}, nil
}

// PerformanceWithScoreAbove retrieve Performance with score above a level
func PerformanceWithScoreAbove(minScore float64) ([]Performance, error) {
	rows, err := db.DB.Query(`SELECT isin, label, score FROM funds WHERE score >= $1 ORDER BY isin`, minScore)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var (
		performances []Performance
		isin         string
		label        string
		score        float64
	)

	for rows.Next() {
		if err := rows.Scan(&isin, &label, &score); err != nil {
			return nil, err
		}

		performances = append(performances, Performance{Isin: isin, Label: label, Score: score})
	}

	return performances, nil
}

// SavePerformance create or update given Performance
func SavePerformance(perf Performance, tx *sql.Tx) (err error) {
	var usedTx *sql.Tx

	if usedTx, err = db.GetTx(tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = PerformanceByIsin(perf.Isin); err != nil {
		err = createPerformance(perf, usedTx)
	} else {
		err = updatePerformance(perf, usedTx)
	}

	return
}

func createPerformance(perf Performance, tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO funds (isin, label, score) VALUES ($1, $2, $3)`, perf.Isin, perf.Label, perf.Score)

	return err
}

func updatePerformance(perf Performance, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE funds SET score = $1, update_date = $2 WHERE isin = $3`, perf.Score, `now()`, perf.Isin)

	return err
}
