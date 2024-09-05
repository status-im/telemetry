package metrics

import (
	"database/sql"

	"github.com/status-im/telemetry/pkg/types"
)

func InsertTelemetryRecord(tx *sql.Tx, data *types.TelemetryRecord) (int, error) {
	result := tx.QueryRow(`
		INSERT INTO telemetryRecord (nodeName, peerId, statusVersion, deviceType)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, data.NodeName, data.PeerID, data.StatusVersion, data.DeviceType)
	if result.Err() != nil {
		return 0, result.Err()
	}

	var recordId int
	err := result.Scan(&recordId)
	if err != nil {
		return 0, err
	}

	return int(recordId), nil
}
