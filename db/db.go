package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Not referenced but needed for database/sql
	_ "github.com/lib/pq"
)

// DB configured or nil if not
var DB *sql.DB

// InitDB start DB connection
func InitDB() {
	dbHost := os.Getenv(`FUNDS_DATABASE_HOST`)
	dbPort := os.Getenv(`FUNDS_DATABASE_PORT`)
	dbUser := os.Getenv(`FUNDS_DATABASE_USER`)
	dbPass := os.Getenv(`FUNDS_DATABASE_PASS`)
	dbName := os.Getenv(`FUNDS_DATABASE_NAME`)

	if dbHost == `` {
		return
	}

	database, err := sql.Open(`postgres`, fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`Connected to %s database`, dbName)
	DB = database
}

// GetTx return given transaction if not nil or create a new one
func GetTx(tx *sql.Tx) (*sql.Tx, error) {
	if tx == nil {
		return DB.Begin()
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
