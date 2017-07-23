package model

import (
	"database/sql"
	"fmt"
	"log"

	// Not referenced but needed for database/sql
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB start DB connection
func InitDB(dbHost string, dbPort int, dbUser string, dbPass string, dbName string) {
	database, err := sql.Open(`postgres`, fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`, dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db = database
}

// RetrieveByID retrieve Performance from database by isin
func RetrieveByID(isin string) (Performance, error) {
	perf := Performance{Isin: isin}

	var score float64
	err := db.QueryRow(`SELECT score FROM funds WHERE isin=$1`, isin).Scan(&score)

	if err != nil {
		return perf, err
	}

	perf.Score = score
	return perf, nil
}

// Save create or update given Performance
func Save(perf Performance, tx *sql.Tx) error {
	var err error
	var usedTx *sql.Tx

	if tx == nil {
		usedTx, err = db.Begin()
		if err != nil {
			return err
		}
	}

	if _, err = RetrieveByID(perf.Isin); err != nil {
		err = create(perf, usedTx)
	} else {
		err = update(perf, usedTx)
	}

	if tx != usedTx {
		if err != nil {
			usedTx.Rollback()
		} else {
			usedTx.Commit()
		}
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
