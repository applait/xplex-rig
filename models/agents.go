package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
)

// Agent model stores basic information of each agent
type Agent struct {
	ID     string `sql:",pk"`
	Secret string `sql:",notnull"`
	Host   string
	Region string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// genSecret generates multistream keys
func genSecret(id string) string {
	t := fmt.Sprintf("xplex://%s_%s@agent", id, uuid.NewV4().String())
	s1 := sha1.New()
	s1.Write([]byte(t))
	return hex.EncodeToString(s1.Sum(nil))
}

// Create attempts creating a new Agent based on given `ID`. `ID` used should be
// hostname for agent.
func (a *Agent) Create(db *pg.DB) (bool, error) {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.Secret = genSecret(a.ID)
	res, err := db.Model(a).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// Update current agent in DB
func (a *Agent) Update(db *pg.DB) (bool, error) {
	a.UpdatedAt = time.Now()
	res, err := db.Model(a).Update()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// UpdateSecret generates a new secret for agent and updates it in DB
func (a *Agent) UpdateSecret(db *pg.DB) (bool, error) {
	a.Secret = genSecret(a.ID)
	return a.Update(db)
}
