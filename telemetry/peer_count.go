package telemetry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/telemetry/pkg/types"
)

type PeerCount struct {
	data types.PeerCount
}

func (r *PeerCount) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer count: %v", err))
		return err
	}

	stmt, err := db.Prepare("INSERT INTO peerCount (timestamp, nodeName, nodeKeyUid, peerId, peerCount, statusVersion, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;")
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
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving peer count: %v", err))
		return err
	}
	r.data.ID = lastInsertId

	return nil
}

type PeerConnFailure struct {
	data types.PeerConnFailure
}

func (r *PeerConnFailure) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer connection failure: %v", err))
		return err
	}

	stmt, err := db.Prepare("INSERT INTO peerConnFailure (timestamp, nodeName, nodeKeyUid, peerId, failedPeerId, failureCount, statusVersion, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
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
		r.data.PeerId,
		r.data.FailedPeerId,
		r.data.FailureCount,
		r.data.StatusVersion,
		r.data.CreatedAt,
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving peer connection failure: %v", err))
		return err
	}
	r.data.ID = lastInsertId

	return nil
}
