package models

import (
	"log"
	"time"

	"github.com/go-pg/pg"
)

// ConnectPG connects to a Postgres database given a connection URI and returns
// the DB interface with logging set up
func ConnectPG(dburi string) (*pg.DB, error) {
	parsed, err := pg.ParseURL(dburi)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(parsed)
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})
	return db, nil
}
