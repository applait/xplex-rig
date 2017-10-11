package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/token"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
)

// MultiStream stores information of multi-streaming keys for users
type MultiStream struct {
	ID            uuid.UUID `sql:",pk,type:uuid"`
	Key           string    `sql:",unique,notnull"`
	IsActive      bool      `sql:",notnull,default:false"`
	IsStreaming   bool      `sql:",notnull,default:false"`
	UserAccountID uuid.UUID `sql:",type:uuid"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	UserAccount *UserAccount // MultiStream belongsTo UserAccount
	Outputs     []*Output    // MultiStream hasMany Output
}

// Output stores configuration information of RTMP ingestion services
// to push to
type Output struct {
	ID       int
	Service  string `sql:",notnull"`
	Key      string `sql:",notnull"`
	Server   string `sql:",notnull,default:'default'"`
	IsActive bool   `sql:",notnull,default:false"`

	MultiStreamID uuid.UUID    `sql:",type:uuid"`
	MultiStream   *MultiStream // Output belongsTo MultiStream
}

// genKey generates multistream keys
func genKey(userid uuid.UUID) (string, error) {
	k := uuid.NewV4().String()
	t, err := token.NewUserToken(userid, fmt.Sprintf("xplex://%s@%s", userid, k), k)
	if err != nil {
		return "", err
	}
	s1 := sha1.New()
	s1.Write([]byte(t))
	return hex.EncodeToString(s1.Sum(nil)), nil
}

// Find a MultiStream
func (m *MultiStream) Find(db *pg.DB) error {
	q := db.Model(m)
	if m.ID != uuid.Nil {
		q.Where("multi_stream.id = ?", m.ID)
	}
	if m.Key != "" {
		q.Where("multi_stream.key = ?", m.Key)
	}
	if m.IsActive {
		q.Where("multi_stream.is_active = ?", true)
	}
	if m.IsStreaming {
		q.Where("multi_stream.is_streaming = ?", true)
	}
	if m.UserAccountID != uuid.Nil {
		q.Where("user_account.id = ?", m.UserAccountID)
		q.Column("UserAccount.id", "UserAccount.username", "UserAccount.is_active")
	}
	return q.First()
}

// Create generates Key and inserts MultiStream row in DB
func (m *MultiStream) Create(db *pg.DB) error {
	u := UserAccount{ID: m.UserAccountID}
	err := u.Find(db)
	if err != nil {
		return err
	}
	if m.Key, err = genKey(m.UserAccountID); err != nil {
		return err
	}
	m.IsActive = true
	return m.Insert(db)
}

// Update current multistream in DB
func (m *MultiStream) Update(db *pg.DB) (bool, error) {
	m.UpdatedAt = time.Now()
	res, err := db.Model(m).Update()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// UpdateKey updates
func (m *MultiStream) UpdateKey(db *pg.DB) (bool, error) {
	key, err := genKey(m.UserAccountID)
	if err != nil {
		return false, err
	}
	m.Key = key
	return m.Update(db)
}

// Insert inserts new row for MultiStream after adding timestamps
func (m *MultiStream) Insert(db *pg.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.ID = uuid.NewV4()
	return db.Insert(m)
}

// Insert inserts new row for Output
func (o *Output) Insert(db *pg.DB) error {
	s, ok := config.MSServices[o.Service]
	if !ok {
		return errors.New("Invalid service name provided")
	}
	if o.Server == "" {
		o.Server = "default"
	}
	if !s.HasServer(o.Server) {
		return errors.New("Invalid server name provided")
	}
	o.IsActive = true
	return db.Insert(o)
}

type resOutputDest struct {
	ConfigID int    `json:"configID"`
	Service  string `json:"service"`
	IsActive bool   `json:"isActive"`
	URL      string `json:"url"`
}

// ResStreamConfig produces output for stream configurations
type ResStreamConfig struct {
	StreamID     uuid.UUID       `json:"streamID"`
	StreamKey    string          `json:"streamKey"`
	IsActive     bool            `json:"isActive"`
	Destinations []resOutputDest `json:"destinations"`
}

// UserStreams returns config of streams for given user
func UserStreams(uid uuid.UUID, isStreaming bool, isActive bool, db *pg.DB) ([]ResStreamConfig, error) {
	var ms []MultiStream
	var o []ResStreamConfig
	err := db.Model(&ms).
		Column("Outputs").
		Where("multi_stream.user_account_id = ?", uid).
		Where("multi_stream.is_streaming = ?", isStreaming).
		Where("multi_stream.is_active = ?", isActive).
		Select()
	if err != nil {
		return nil, err
	}
	for _, m := range ms {
		el := ResStreamConfig{
			StreamID:  m.ID,
			StreamKey: m.Key,
			IsActive:  m.IsActive,
		}
		for _, mo := range m.Outputs {
			s := config.MSServices[mo.Service]
			odest := resOutputDest{
				ConfigID: mo.ID,
				Service:  mo.Service,
				IsActive: mo.IsActive,
				URL:      s.RTMPUrl(mo.Key, mo.Server),
			}
			el.Destinations = append(el.Destinations, odest)
		}
		o = append(o, el)
	}
	return o, nil
}
