package metrics

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type ReceivedEnvelope struct {
	types.ReceivedEnvelope
}

func (r *ReceivedEnvelope) put(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO receivedEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, receiverKeyUID, nodeName, processingError, statusVersion, deviceType)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(
		r.MessageHash,
		r.SentAt,
		r.CreatedAt,
		r.PubsubTopic,
		r.Topic,
		r.ReceiverKeyUID,
		r.NodeName,
		r.ProcessingError,
		r.StatusVersion,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}
	r.ID = lastInsertId

	return nil
}

func (r *ReceivedEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding received envelope: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving received envelope: %v", err))
		return err
	}

	return nil
}

func (r *ReceivedEnvelope) UpdateProcessingError(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`UPDATE receivedEnvelopes SET processingError=$1 WHERE
							messageHash = $2 AND sentAt = $3 AND
							pubsubTopic = $4 AND topic = $5 AND
							receiverKeyUID = $6 AND nodeName = $7;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.ProcessingError, r.MessageHash, r.SentAt, r.PubsubTopic, r.Topic, r.ReceiverKeyUID, r.NodeName)
	if err != nil {
		return err
	}

	return nil
}

type SentEnvelope struct {
	types.SentEnvelope
}

func (r *SentEnvelope) put(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO sentEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion, deviceType)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
							ON CONFLICT ON CONSTRAINT sentEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		r.MessageHash,
		r.SentAt,
		r.CreatedAt,
		r.PubsubTopic,
		r.Topic,
		r.SenderKeyUID,
		r.PeerID,
		r.NodeName,
		r.PublishMethod,
		r.StatusVersion,
		r.DeviceType,
	).Scan(&lastInsertId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	defer stmt.Close()
	r.ID = int(lastInsertId)

	return nil
}
func (r *SentEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
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
	types.ErrorSendingEnvelope
}

func (e *ErrorSendingEnvelope) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &e); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding error sending envelope: %v", err))
		return err
	}

	e.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO errorSendingEnvelope (messageHash, sentAt, createdAt, pubsubTopic,
		topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion, error, deviceType)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT ON CONSTRAINT errorSendingEnvelope_unique DO NOTHING
		RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(
		e.SentEnvelope.MessageHash,
		e.SentEnvelope.SentAt,
		e.CreatedAt,
		e.SentEnvelope.PubsubTopic,
		e.SentEnvelope.Topic,
		e.SentEnvelope.SenderKeyUID,
		e.SentEnvelope.PeerID,
		e.SentEnvelope.NodeName,
		e.SentEnvelope.PublishMethod,
		e.SentEnvelope.StatusVersion,
		e.Error,
		e.DeviceType,
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
	e.SentEnvelope.ID = int(lastInsertId)

	return nil
}
