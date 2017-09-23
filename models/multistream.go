package models

import "time"

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
