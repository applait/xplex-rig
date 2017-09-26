package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/applait/xplex-rig/token"
	"github.com/go-pg/pg"
)

// MultiStream stores information of multi-streaming keys for users
type MultiStream struct {
	ID          int
	Key         string `sql:",unique,notnull"`
	IsActive    bool   `sql:",notnull,default:false"`
	IsStreaming bool   `sql:",notnull,default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	UserID             int                  // MultiStream belongsTo User
	MultiStreamConfigs []*MultiStreamConfig // MultiStream hasMany MultiStreamConfigs
}

func genKey(u *User) (string, error) {
	t, err := token.NewUserToken(u.ID, fmt.Sprintf("xplex://%s@%s", u.Username, time.Now()), u.Password)
	if err != nil {
		return "", err
	}
	s1 := sha1.New()
	s1.Write([]byte(t))
	return hex.EncodeToString(s1.Sum(nil)), nil
}

// Create generates Key and inserts MultiStream row in DB
func (m *MultiStream) Create(db *pg.DB) error {
	u := User{ID: m.UserID}
	err := u.Find(db)
	if err != nil {
		return err
	}
	if m.Key, err = genKey(&u); err != nil {
		return err
	}
	return m.Insert(db)
}

// Insert inserts new row for MultiStream after adding timestamps
func (m *MultiStream) Insert(db *pg.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return db.Insert(m)
}

// MultiStreamConfig stores configuration information of RTMP ingestion services
// to push to
type MultiStreamConfig struct {
	ID       int
	Service  string `sql:",notnull"`
	Key      string `sql:",notnull"`
	Server   string `sql:",notnull,default:'default'"`
	IsActive bool   `sql:",notnull,default:false"`

	MultiStreamID int // MultiStreamConfig belongsTo MultiStream
}
