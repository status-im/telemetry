package metrics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type PeerCount struct {
	data types.PeerCount
}

func (r *PeerCount) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer count: %v", err))
		return err
	}

	stmt, err := db.Prepare("INSERT INTO peerCount (timestamp, nodeName, nodeKeyUid, peerId, peerCount, statusVersion, createdAt, deviceType) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.data.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(
		r.data.Timestamp,
		r.data.NodeName,
		r.data.NodeKeyUid,
		r.data.PeerID,
		r.data.PeerCount,
		r.data.StatusVersion,
		r.data.CreatedAt,
		r.data.DeviceType,
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving peer count: %v", err))
		return err
	}
	r.data.ID = lastInsertId

	return nil
}

type PeerConnFailure struct {
	types.PeerConnFailure
}

func (r *PeerConnFailure) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer connection failure: %v", err))
		return err
	}

	stmt, err := db.Prepare("INSERT INTO peerConnFailure (timestamp, nodeName, nodeKeyUid, peerId, failedPeerId, failureCount, statusVersion, createdAt, deviceType) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(
		r.Timestamp,
		r.NodeName,
		r.NodeKeyUid,
		r.PeerId,
		r.FailedPeerId,
		r.FailureCount,
		r.StatusVersion,
		r.CreatedAt,
		r.DeviceType,
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving peer connection failure: %v", err))
		return err
	}
	r.ID = lastInsertId

	return nil
}
