package telemetry

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/status-im/telemetry/pkg/types"
)

type ReceivedEnvelope struct {
	data types.ReceivedEnvelope
}

func (r *ReceivedEnvelope) put(db *sql.DB) error {
	r.data.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO receivedEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, receiverKeyUID, nodeName, processingError, statusVersion)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(
		r.data.MessageHash,
		r.data.SentAt,
		r.data.CreatedAt,
		r.data.PubsubTopic,
		r.data.Topic,
		r.data.ReceiverKeyUID,
		r.data.NodeName,
		r.data.ProcessingError,
		r.data.StatusVersion,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}
	r.data.ID = lastInsertId

	return nil
}

func (r *ReceivedEnvelope) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding received envelope: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving received envelope: %v", err))
		return err
	}

	return nil
}

func (r *ReceivedEnvelope) updateProcessingError(db *sql.DB) error {
	r.data.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`UPDATE receivedEnvelopes SET processingError=$1 WHERE
							messageHash = $2 AND sentAt = $3 AND
							pubsubTopic = $4 AND topic = $5 AND
							receiverKeyUID = $6 AND nodeName = $7;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.data.ProcessingError, r.data.MessageHash, r.data.SentAt, r.data.PubsubTopic, r.data.Topic, r.data.ReceiverKeyUID, r.data.NodeName)
	if err != nil {
		return err
	}

	return nil
}

type SentEnvelope struct {
	data types.SentEnvelope
}

func (r *SentEnvelope) put(db *sql.DB) error {
	r.data.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO sentEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
							ON CONFLICT ON CONSTRAINT sentEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		r.data.MessageHash,
		r.data.SentAt,
		r.data.CreatedAt,
		r.data.PubsubTopic,
		r.data.Topic,
		r.data.SenderKeyUID,
		r.data.PeerID,
		r.data.NodeName,
		r.data.PublishMethod,
		r.data.StatusVersion,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	defer stmt.Close()
	r.data.ID = int(lastInsertId)

	return nil
}
func (r *SentEnvelope) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding sent envelope: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving sent envelope: %v", err))
		return err
	}
	return nil
}

type ErrorSendingEnvelope struct {
	data types.ErrorSendingEnvelope
}

func (e *ErrorSendingEnvelope) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &e); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding error sending envelope: %v", err))
		return err
	}

	e.data.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO errorSendingEnvelope (messageHash, sentAt, createdAt, pubsubTopic,
		topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT ON CONSTRAINT errorSendingEnvelope_unique DO NOTHING
		RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		e.data.SentEnvelope.MessageHash,
		e.data.SentEnvelope.SentAt,
		e.data.CreatedAt,
		e.data.SentEnvelope.PubsubTopic,
		e.data.SentEnvelope.Topic,
		e.data.SentEnvelope.SenderKeyUID,
		e.data.SentEnvelope.PeerID,
		e.data.SentEnvelope.NodeName,
		e.data.SentEnvelope.PublishMethod,
		e.data.SentEnvelope.StatusVersion,
		e.data.Error,
	).Scan(&lastInsertId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			errs.Append(data.Id, fmt.Sprintf("Error saving error sending envelope: %v", err))
			return err
		}
	}

	defer stmt.Close()
	e.data.SentEnvelope.ID = int(lastInsertId)

	return nil
}
