package telemetry

import (
	"database/sql"

	"github.com/status-im/telemetry/pkg/types"
)

func InsertCommonFields(db *sql.DB, data types.CommonFieldsAccessor) (int, error) {
	stmt, err := db.Prepare(`
		INSERT INTO commonFields (nodeName, peerId, statusVersion, deviceType)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var commonFieldsId int
	err = stmt.QueryRow(
		data.GetNodeName(),
		data.GetPeerID(),
		data.GetStatusVersion(),
		data.GetDeviceType(),
	).Scan(&commonFieldsId)
	if err != nil {
		return 0, err
	}

	return commonFieldsId, nil
}
