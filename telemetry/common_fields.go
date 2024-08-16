package telemetry

import (
	"database/sql"

	"github.com/status-im/telemetry/pkg/types"
)

func InsertCommonFields(tx *sql.Tx, data *types.CommonFields) (int, error) {
	stmt, err := tx.Prepare(`
		INSERT INTO commonFields (nodeName, peerId, statusVersion, deviceType)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`)
	if err != nil {
		return 0, err
	}

	var commonFieldsId int
	err = stmt.QueryRow(
		data.NodeName,
		data.PeerID,
		data.StatusVersion,
		data.DeviceType,
	).Scan(&commonFieldsId)
	if err != nil {
		return 0, err
	}

	return commonFieldsId, nil
}
