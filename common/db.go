package common

import (
	"database/sql"
	"log"

	// Postgresql driver for database/sql
	_ "github.com/lib/pq"
)

// ConnectDB establishes database connection with postgres
func ConnectDB(connStr string) error {
	log.Println("Connecting to DB")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	log.Println("Connected to db")
	return nil
}
