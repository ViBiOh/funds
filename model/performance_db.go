package model

import (
	"database/sql"

	"github.com/ViBiOh/funds/db"
)

// RetrieveByID retrieve Performance from database by isin
func RetrieveByID(isin string) (*Performance, error) {
	var score float64
	err := db.DB.QueryRow(`SELECT score FROM funds WHERE isin=$1`, isin).Scan(&score)

	if err != nil {
		return nil, err
	}

	return &Performance{Isin: isin, Score: score}, nil
}

// SaveAll create or update all given Performances
func SaveAll(performances []Performance, tx *sql.Tx) error {
	var err error
	var usedTx *sql.Tx

	defer func() {
		deferTx(tx, usedTx, err)
	}()

	if usedTx, err = getTx(tx); err != nil {
		return err
	}

	for _, performance := range performances {
		if err = Save(performance, usedTx); err != nil {
			return err
		}
	}

	return nil
}

// Save create or update given Performance
func Save(perf Performance, tx *sql.Tx) error {
	var err error
	var usedTx *sql.Tx

	defer func() {
		deferTx(tx, usedTx, err)
	}()

	if usedTx, err = getTx(tx); err != nil {
		return err
	}

	if _, err = RetrieveByID(perf.Isin); err != nil {
		err = create(perf, usedTx)
	} else {
		err = update(perf, usedTx)
	}

	return err
}

func create(perf Performance, tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO funds (isin, score) VALUES ($1, $2)`, perf.Isin, perf.Score)

	return err
}

func update(perf Performance, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE funds SET score=$1, update_date=$2 WHERE isin=$3`, perf.Score, `now()`, perf.Isin)

	return err
}

func getTx(tx *sql.Tx) (*sql.Tx, error) {
	if tx == nil {
		return db.DB.Begin()
	}

	return tx, nil
}

func deferTx(tx *sql.Tx, usedTx *sql.Tx, err error) {
	if usedTx != tx {
		if err != nil {
			usedTx.Rollback()
		} else {
			usedTx.Commit()
		}
	}
}
