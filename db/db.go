package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Not referenced but needed for database/sql
	_ "github.com/lib/pq"
)

var db *sql.DB

// Init start DB connection
func Init() error {
	dbHost := os.Getenv(`FUNDS_DATABASE_HOST`)
	dbPort := os.Getenv(`FUNDS_DATABASE_PORT`)
	dbUser := os.Getenv(`FUNDS_DATABASE_USER`)
	dbPass := os.Getenv(`FUNDS_DATABASE_PASS`)
	dbName := os.Getenv(`FUNDS_DATABASE_NAME`)

	if dbHost == `` {
		return nil
	}

	database, err := sql.Open(`postgres`, fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		return fmt.Errorf(`Error while opening database connection: %v`, err)
	}

	err = database.Ping()
	if err != nil {
		return fmt.Errorf(`Error while pinging database: %v`, err)
	}

	db = database

	return nil
}

// Ping indicate if database is ready or not
func Ping() bool {
	return db != nil && db.Ping() == nil
}

// GetTx return given transaction if not nil or create a new one
func GetTx(label string, tx *sql.Tx) (*sql.Tx, error) {
	if tx == nil {
		usedTx, err := db.Begin()

		if err != nil {
			return nil, fmt.Errorf(`Error while getting transaction for %s: %v`, label, err)
		}
		return usedTx, nil
	}

	return tx, nil
}

// EndTx ends transaction according error without shadowing given error
func EndTx(label string, tx *sql.Tx, err error) error {
	if err != nil {
		if endErr := tx.Rollback(); endErr != nil {
			log.Printf(`Error while rolling back transaction for %s: %v`, label, endErr)
		}
	} else if endErr := tx.Commit(); endErr != nil {
		return fmt.Errorf(`Error while committing transaction for %s: %v`, label, endErr)
	}

	return nil
}

// RowsClose closes rows without shadowing error
func RowsClose(label string, rows *sql.Rows, err error) error {
	if endErr := rows.Close(); endErr != nil {
		endErr = fmt.Errorf(`Error while closing rows for %s: %v`, label, endErr)

		if err == nil {
			return endErr
		}
		log.Print(endErr)
	}

	return err
}

// Query wraps https://golang.org/pkg/database/sql/#DB.Query
func Query(label string, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return rows, fmt.Errorf(`Error while querying %s: %v`, label, err)
	}
	return rows, err
}

// QueryRow wraps https://golang.org/pkg/database/sql/#DB.QueryRow
func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}
