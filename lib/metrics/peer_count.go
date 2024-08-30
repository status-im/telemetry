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

func (r *PeerCount) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer count: %v", err))
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.data.TelemetryRecord)
	if err != nil {
		return err
	}

	peerCountStmt, err := tx.Prepare(`
		INSERT INTO peerCount (common_fields_id, nodeKeyUid, peerCount)
		VALUES ($1, $2, $3)
		RETURNING id;
	`)
	if err != nil {
		return err
	}

	var lastInsertId int
	err = peerCountStmt.QueryRow(
		recordId,
		r.data.NodeKeyUID,
		r.data.PeerCount,
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving peer count: %v", err))
		return err
	}
	r.ID = lastInsertId

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

func (r *PeerConnFailure) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding peer connection failure: %v", err))
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.data.TelemetryRecord)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO peerConnFailure (common_fields_id, nodeKeyUid, failedPeerId, failureCount) VALUES ($1, $2, $3, $4) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	lastInsertId := 0
	err = stmt.QueryRow(
		recordId,
		r.data.NodeKeyUID,
		r.data.FailedPeerId,
		r.data.FailureCount,
	).Scan(&lastInsertId)
	if err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving peer connection failure: %v", err))
		return err
	}
	r.ID = lastInsertId

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

func (r *PeerConnFailure) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "peerConnFailure", before)
}
