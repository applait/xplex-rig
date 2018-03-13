package stream

import (
	"database/sql"
	"fmt"

	"github.com/applait/xplex-rig/common"
	uuid "github.com/satori/go.uuid"
)

// GetDestinations fetches all destinations created for a stream
func GetDestinations(streamID uuid.UUID) ([]common.Destination, error) {
	var ds []common.Destination
	query := `
    select
      id, service, stream_key, is_active
    from destinations
      where multi_stream_id = $1;
  `
	rows, err := common.DB.Query(query, streamID)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return ds, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch destinations for stream")
	}
	for rows.Next() {
		var d common.Destination
		err = rows.Scan(
			&d.ID, &d.Service, &d.StreamKey, &d.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("Unable to fetch destinations for stream")
		}
		ds = append(ds, d)
	}
	err = rows.Err()
	if err != nil {
		return ds, err
	}
	return ds, nil
}

// AddDestination adds a new stream destination
func AddDestination(streamID uuid.UUID, service string, key string, userID uuid.UUID) (common.Destination, error) {
	var d common.Destination
	s, ok := Services[service]
	if !ok {
		return d, fmt.Errorf("Invalid service name: %s", service)
	}
	d.Service = service
	d.StreamKey = key
	d.IsActive = true
	query := `
    insert into destinations
      (service, stream_key, is_active, multi_stream_id)
    values
      ($1, $2, $3, $4)
    returning id;
  `
	err := common.DB.QueryRow(query, d.Service, d.StreamKey, true, streamID).Scan(&d.ID)
	if err != nil {
		return d, fmt.Errorf("Unable to add stream destination for %s", service)
	}
	d.RTMPUrl = s.RTMPUrl(d.StreamKey, "default")
	return d, nil
}
