package db

import (
	"github.com/jinzhu/gorm"

	// Postgres dialet for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB holds the DB connection pool returned by GORM
type DB struct {
	DB  *gorm.DB
	URI string
}

// Close closes the underlying database connection
func (db *DB) Close() {
	db.DB.Close()
}

// Migrate auto-migrates given models
func (db *DB) migrate(m interface{}) {
	db.DB.AutoMigrate(m)
}

// NewDB connects to a Postgres database given a connection URI and returns `DB`
// if connection is established.
func NewDB(dburi string) (*DB, error) {
	db, err := gorm.Open(dburi)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db, URI: dburi}, nil
}
