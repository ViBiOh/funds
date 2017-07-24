package db

import (
	"database/sql"
	"fmt"
	"log"

	// Not referenced but needed for database/sql
	_ "github.com/lib/pq"
)

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

// DeferTx check if transaction is hold by current function and end it properly
func DeferTx(tx *sql.Tx, usedTx *sql.Tx, err error) {
	if usedTx != tx {
		if err != nil {
			usedTx.Rollback()
		} else {
			usedTx.Commit()
		}
	}
}
