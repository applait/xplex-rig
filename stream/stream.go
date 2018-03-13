package stream

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/applait/xplex-rig/common"
	"github.com/satori/go.uuid"
)

// createStreamKey generates a base64 encoded UUID v5 for streams given user ID and stream ID
func createStreamKey(userID uuid.UUID, streamID uuid.UUID) string {
	uid := uuid.NewV5(userID, fmt.Sprintf("%s-%s", streamID.String(), time.Now()))
	return base64.RawURLEncoding.EncodeToString(uid.Bytes())
}

// GetStreamByID returns a stream's information given its ID
func GetStreamByID(streamID uuid.UUID) (common.Stream, error) {
	query := `
    select
      m.id, m.stream_key, m.is_active, m.is_streaming, m.user_account_id
    from multi_streams m
      where m.id = $1;
  `
	var s common.Stream
	err := common.DB.QueryRow(query, streamID).Scan(
		&s.ID, &s.StreamKey, &s.IsActive, &s.IsStreaming, &s.User.ID,
	)
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetStreamByStreamKey returns a stream's information given a stream key
func GetStreamByStreamKey(streamKey string) (common.Stream, error) {
	query := `
    select
      m.id, m.stream_key, m.is_active, m.is_streaming, m.user_account_id
    from multi_streams m
      where m.stream_key = $1;
  `
	var s common.Stream
	err := common.DB.QueryRow(query, streamKey).Scan(
		&s.ID, &s.StreamKey, &s.IsActive, &s.IsStreaming, &s.User.ID,
	)
	if err != nil {
		return s, err
	}
	return s, nil
}

// CreateStream creates a new stream for a given user ID
func CreateStream(userID uuid.UUID) (common.Stream, error) {
	query := `
    insert into multi_streams
      (id, stream_key, is_active, is_streaming, user_account_id, created_at, updated_at)
    values
      ($1, $2, $3, $4, $5, now(), now());
  `
	var s common.Stream
	s.ID = uuid.NewV4()
	s.StreamKey = createStreamKey(userID, s.ID)
	s.IsActive = true
	s.IsStreaming = false
	_, err := common.DB.Exec(query,
		&s.ID,
		&s.StreamKey,
		&s.IsActive,
		&s.IsStreaming,
		&userID,
	)
	if err != nil {
		return s, err
	}
	return s, nil
}

// ChangeStreamKey creates a new stream key for a given stream of a given user
func ChangeStreamKey(streamID uuid.UUID, userID uuid.UUID) (string, error) {
	query := `
    update multi_streams
      set stream_key = $1
    where id = $2;
  `
	newKey := createStreamKey(userID, streamID)
	res, err := common.DB.Exec(query, newKey, streamID)
	if err != nil {
		return newKey, err
	}
	if affected, err := res.RowsAffected(); err != nil || affected != 1 {
		return newKey, errors.New("Unable to update stream key")
	}
	return newKey, nil
}
