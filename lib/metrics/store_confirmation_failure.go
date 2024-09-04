package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type StoreConfirmationFailed struct {
	types.StoreConfirmationFailed
}

func (s *StoreConfirmationFailed) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &s); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding store confirmation failed: %v", err))
		return err
	}

	if err := s.Put(ctx, db); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving store confirmation failed: %v", err))
		return err
	}

	log.Printf("AK: store confirmation metric saved: %v", s)
	return nil
}

func (s *StoreConfirmationFailed) Put(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &s.TelemetryRecord)
	if err != nil {
		return err
	}
	result := tx.QueryRow("INSERT INTO storeConfirmationFailed (recordId, messageHash) VALUES ($1, $2) RETURNING id;", recordId, s.MessageHash)
	if result.Err() != nil {
		return result.Err()
	}

	err = result.Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreConfirmationFailed) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "storeConfirmationFailed", before)
}
