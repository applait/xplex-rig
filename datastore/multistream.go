package datastore

import "time"

// RTMPStream stores information of multi-streaming keys for users
type RTMPStream struct {
	ID          int
	Key         string `sql:",unique,notnull"`
	IsActive    bool   `sql:",notnull,default:false"`
	IsStreaming bool   `sql:",notnull,default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	UserID     int         // RTMPStream belongsTo User
	RTMPPushes []*RTMPPush // RTMPStream hasMany RTMPPush
}

// RTMPPush store configuration information of RTMP ingestion services
// to push to
type RTMPPush struct {
	ID       int
	Service  string `sql:",notnull"`
	Key      string `sql:",notnull"`
	Server   string `sql:",notnull,default:'default'"`
	IsActive bool   `sql:",notnull,default:false"`

	RTMPStreamID int // RTMPPush belongsTo RTMPStream
}
