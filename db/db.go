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
		return fmt.Errorf(`Error while opening: %v`, err)
	}

	err = database.Ping()
	if err != nil {
		return fmt.Errorf(`Error while pinging: %v`, err)
	}

	db = database

	return nil
}

// Ping indicate if database is ready or not
func Ping() bool {
	return db != nil && db.Ping() == nil
}

// GetTx return given transaction if not nil or create a new one
func GetTx(tx *sql.Tx) (*sql.Tx, error) {
	if tx == nil {
		return db.Begin()
	}

	return tx, nil
}

// EndTx end transaction properly according to error
func EndTx(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf(`Error while rolling back transaction: %v`, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

// Query wraps https://golang.org/pkg/database/sql/#DB.Query
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// QueryRow wraps https://golang.org/pkg/database/sql/#DB.QueryRow
func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}
