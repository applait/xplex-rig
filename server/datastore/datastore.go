package datastore

import (
	"log"
	"time"

	"github.com/go-pg/pg"
)

// DataStore defines a basic common interface common to all data store operations
type DataStore interface {
	Insert(*pg.DB) error
	Update(*pg.DB) error
}

// PGModel provides common fields and data for postgres models to add useful methods
type PGModel struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
