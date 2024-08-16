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

func (r *ReceivedEnvelope) put(db *sql.DB) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	commonFieldsId, err := InsertCommonFields(tx, &r.data.CommonFields)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO receivedEnvelopes (commonFieldsId, messageHash, sentAt, pubsubTopic,
							topic, receiverKeyUID, processingError)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(
		commonFieldsId,
		r.data.MessageHash,
		r.data.SentAt,
		r.data.PubsubTopic,
		r.data.Topic,
		r.data.ReceiverKeyUID,
		r.data.ProcessingError,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}
	r.ID = lastInsertId

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *ReceivedEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding received envelope: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving received envelope: %v", err))
		return err
	}

	return nil
}

func (r *ReceivedEnvelope) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "receivedEnvelopes", before)
}

func (r *ReceivedEnvelope) UpdateProcessingError(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
func (r *ReceivedEnvelope) updateProcessingError(db *sql.DB) error {
	stmt, err := db.Prepare(`UPDATE receivedEnvelopes SET processingError=$1 WHERE
							messageHash = $2 AND sentAt = $3 AND
							pubsubTopic = $4 AND topic = $5 AND
							receiverKeyUID = $6;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.data.ProcessingError, r.data.MessageHash, r.data.SentAt, r.data.PubsubTopic, r.data.Topic, r.data.ReceiverKeyUID)
	if err != nil {
		return err
	}

	return nil
}

type SentEnvelope struct {
	types.SentEnvelope
}

func (r *SentEnvelope) put(db *sql.DB) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	commonFieldsId, err := InsertCommonFields(tx, &r.data.CommonFields)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO sentEnvelopes (commonFieldsId, messageHash, sentAt, pubsubTopic,
							topic, senderKeyUID, publishMethod)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT ON CONSTRAINT sentEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		commonFieldsId,
		r.data.MessageHash,
		r.data.SentAt,
		r.data.PubsubTopic,
		r.data.Topic,
		r.data.SenderKeyUID,
		r.data.PublishMethod,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	r.data.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
func (r *SentEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding sent envelope: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
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

func (e *ErrorSendingEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &e); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding error sending envelope: %v", err))
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	commonFieldsId, err := InsertCommonFields(tx, &e.data.SentEnvelope.CommonFields)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO errorSendingEnvelope (commonFieldsId, messageHash, sentAt, pubsubTopic,
		topic, senderKeyUID, publishMethod, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ON CONSTRAINT errorSendingEnvelope_unique DO NOTHING
		RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		commonFieldsId,
		e.data.SentEnvelope.MessageHash,
		e.data.SentEnvelope.SentAt,
		e.data.SentEnvelope.PubsubTopic,
		e.data.SentEnvelope.Topic,
		e.data.SentEnvelope.SenderKeyUID,
		e.data.SentEnvelope.PublishMethod,
		e.data.Error,
	).Scan(&lastInsertId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			errs.Append(data.ID, fmt.Sprintf("Error saving error sending envelope: %v", err))
			return err
		}
	}

	e.data.SentEnvelope.ID = int(lastInsertId)

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

func (r *ErrorSendingEnvelope) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "errorSendingEnvelope", before)
}
