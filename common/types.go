// Shared types exported by the common package

package common

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

// DB holds the connection to Postgresql database
var DB *sql.DB

// Config holds instance of parsed config
var Config JSONConfig

// JWTKeys defines the different keys used by rig to sign and verify JWTs for different use cases. These keys need to be
// shared across all rig instances.
type JWTKeys struct {
	Users  string `json:"users"`
	Agents string `json:"agents"`
	Admins string `json:"admins"`
}

// JSONConfig holds the structure for configuration
type JSONConfig struct {
	Port        int     `json:"port"`
	JWTKeys     JWTKeys `json:"jwtKeys"`
	PostgresURL string  `json:"postgresUrl"`
}

// UserAccount represents a single user account
type UserAccount struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	IsActive  bool      `json:"isActive,omitempty"`
}

// Destination represents a single egress destination for Stream
type Destination struct {
	ID        int64  `json:"id"`
	Service   string `json:"service"`
	StreamKey string `json:"streamKey"`
	RTMPUrl   string `json:"rtmpUrl"`
	IsActive  bool   `json:"isActive,omitempty"`
}

// Region represents a region where edge servers are hosted and streams are provisioned
type Region struct {
	DNS           string `json:"dns"`
	Name          string `json:"name"`
	City          string `json:"city"`
	Country       string `json:"country"`
	InfraProvider string `json:"infraProvider,omitempty"`
}

// EdgeServer represents a single edge server that ingests RTMP streams and pushes to Destination configured by a Stream
type EdgeServer struct {
	Hostname  string    `json:"hostname"`
	Region    Region    `json:"region"`
	LastCheck time.Time `json:"lastCheck"`
}

// StreamProvision represents provisioning information of a stream
type StreamProvision struct {
	Region    Region     `json:"region"`
	Server    EdgeServer `json:"server"`
	StreamURL string     `json:"streamUrl"`
}

// Stream represents a single multi-streaming configuration
type Stream struct {
	ID              uuid.UUID       `json:"id"`
	StreamKey       string          `json:"streamKey"`
	Destinations    []Destination   `json:"destinations"`
	ProvisionStatus StreamProvision `json:"provisionStatus"`
	IsStreaming     bool            `json:"isStreaming"`
	User            UserAccount     `json:"user,omitempty"`
	IsActive        bool            `json:"isActive,omitempty"`
}
