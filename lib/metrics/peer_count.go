package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type PeerCount struct {
	types.PeerCount
}

func (r *PeerCount) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding peer count: %v", err))
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.TelemetryRecord)
	if err != nil {
		return err
	}

	result := tx.QueryRow(`
		INSERT INTO peerCount (recordId, nodeKeyUid, peerCount, timestamp)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, recordId, r.NodeKeyUID, r.PeerCount.PeerCount, r.PeerCount.Timestamp)
	if result.Err() != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving peer count: %v", result.Err()))
		return result.Err()
	}

	var lastInsertId int
	err = result.Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

func (r *PeerCount) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "peerCount", before)
}

type PeerConnFailure struct {
	types.PeerConnFailure
}

func (r *PeerConnFailure) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding peer connection failure: %v", err))
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.TelemetryRecord)
	if err != nil {
		return err
	}

	result := tx.QueryRow(`
		INSERT INTO peerConnFailure (recordId, nodeKeyUid, failedPeerId, failureCount, timestamp)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`, recordId, r.NodeKeyUID, r.FailedPeerId, r.FailureCount, r.Timestamp)
	if result.Err() != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving peer connection failure: %v", result.Err()))
		return result.Err()
	}

	var lastInsertId int
	err = result.Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

func (r *PeerConnFailure) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "peerConnFailure", before)
}
