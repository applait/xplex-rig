package stream

import (
	"encoding/base64"
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
