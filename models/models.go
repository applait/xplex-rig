package models

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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

// CreateSchema creates tables in Postgres from models
func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		&User{},
		&MultiStream{},
		&MultiStreamConfig{},
		&EdgeCluster{},
		&AgentNode{},
	}
	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
