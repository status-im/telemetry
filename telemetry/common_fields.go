package telemetry

import (
	"database/sql"

	"github.com/status-im/telemetry/pkg/types"
)

func InsertTelemetryRecord(tx *sql.Tx, data *types.TelemetryRecord) (int, error) {
	stmt, err := tx.Prepare(`
		INSERT INTO telemetryRecord (nodeName, peerId, statusVersion, deviceType)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`)
	if err != nil {
		return 0, err
	}

	var recordId int
	err = stmt.QueryRow(
		data.NodeName,
		data.PeerID,
		data.StatusVersion,
		data.DeviceType,
	).Scan(&recordId)
	if err != nil {
		return 0, err
	}

	return recordId, nil
}
