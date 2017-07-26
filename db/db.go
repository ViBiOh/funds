package db

import (
	"database/sql"
	"fmt"
	"log"

	// Not referenced but needed for database/sql
	_ "github.com/lib/pq"
)

// CommitStep size for batch update processing
const CommitStep = 100

// DB configured or nil if not
var DB *sql.DB

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
		return tx.Rollback()
	}

	return tx.Commit()
}
