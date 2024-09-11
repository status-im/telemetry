package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type ReceivedEnvelope struct {
	types.ReceivedEnvelope
}

func (r *ReceivedEnvelope) put(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.TelemetryRecord)
	if err != nil {
		return err
	}

	result := tx.QueryRow(`INSERT INTO receivedEnvelopes (recordId, messageHash, sentAt, pubsubTopic,
							topic, receiverKeyUID, processingError)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`, recordId, r.MessageHash, r.SentAt, r.PubsubTopic, r.Topic, r.ReceiverKeyUID, r.ProcessingError)
	if result.Err() != nil {
		if errors.Is(result.Err(), sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	var lastInsertId int
	err = result.Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *ReceivedEnvelope) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding received envelope: %v", err))
		return err
	}

	if err := r.put(ctx, db); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving received envelope: %v", err))
		return err
	}

	return nil
}

func (r *ReceivedEnvelope) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "receivedEnvelopes", before)
}

func (r *ReceivedEnvelope) UpdateProcessingError(db *sql.DB) error {
	_, err := db.Exec(`UPDATE receivedEnvelopes SET processingError=$1 WHERE
							messageHash = $2 AND sentAt = $3 AND
							pubsubTopic = $4 AND topic = $5 AND
							receiverKeyUID = $6;`, r.ProcessingError, r.MessageHash, r.SentAt, r.PubsubTopic, r.Topic, r.ReceiverKeyUID)
	if err != nil {
		return err
	}

	return nil
}

type SentEnvelope struct {
	types.SentEnvelope
}

func (r *SentEnvelope) put(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.TelemetryRecord)
	if err != nil {
		return err
	}

	result := tx.QueryRow(`INSERT INTO sentEnvelopes (recordId, messageHash, sentAt, pubsubTopic,
							topic, senderKeyUID, publishMethod)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT ON CONSTRAINT sentEnvelopes_unique DO NOTHING
							RETURNING id;`, recordId, r.MessageHash, r.SentAt, r.PubsubTopic, r.Topic, r.SenderKeyUID, r.PublishMethod)
	if result.Err() != nil {
		if errors.Is(result.Err(), sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	var lastInsertId int
	err = result.Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *SentEnvelope) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding sent envelope: %v", err))
		return err
	}

	if err := r.put(ctx, db); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving sent envelope: %v", err))
		return err
	}
	return nil
}

func (r *SentEnvelope) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "sentEnvelopes", before)
}

type ErrorSendingEnvelope struct {
	types.ErrorSendingEnvelope
}

func (e *ErrorSendingEnvelope) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &e); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding error sending envelope: %v", err))
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &e.SentEnvelope.TelemetryRecord)
	if err != nil {
		return err
	}

	result := tx.QueryRow(`INSERT INTO errorSendingEnvelope (recordId, messageHash, sentAt, pubsubTopic,
		topic, senderKeyUID, publishMethod, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ON CONSTRAINT errorSendingEnvelope_unique DO NOTHING
		RETURNING id;`, recordId, e.SentEnvelope.MessageHash, e.SentEnvelope.SentAt, e.SentEnvelope.PubsubTopic, e.SentEnvelope.Topic, e.SentEnvelope.SenderKeyUID, e.SentEnvelope.PublishMethod, e.Error)

	if result.Err() != nil {
		if errors.Is(result.Err(), sql.ErrNoRows) {
			return nil
		} else {
			errs.Append(data.ID, fmt.Sprintf("Error saving error sending envelope: %v", result.Err()))
			return result.Err()
		}
	}

	var lastInsertId int
	err = result.Scan(&lastInsertId)
	if err != nil {
		return err
	}
	e.SentEnvelope.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

func (r *ErrorSendingEnvelope) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "errorSendingEnvelope", before)
}
