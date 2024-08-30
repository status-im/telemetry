package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type ReceivedMessageAggregated struct {
	ID                int
	ChatID            string
	DurationInSeconds int64
	Value             float64
	RunAt             int64
}

type ReceivedMessage struct {
	types.ReceivedMessage
}

func (r *ReceivedMessage) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding received message failure: %v", err))
		return err
	}

	if err := r.Put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving received messages: %v", err))
		return err
	}
	return nil
}

func (r *ReceivedMessage) Clean(db *sql.DB, before int64) (int64, error) {
	return common.Cleanup(db, "receivedMessages", before)
}

func (r *ReceivedMessage) Put(db *sql.DB) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &r.data.TelemetryRecord)
	if err != nil {
		return fmt.Errorf("failed to insert common fields: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO receivedMessages (recordId, chatId, messageHash, messageId, receiverKeyUID, sentAt, topic, messageType, messageSize, pubSubTopic) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	lastInsertId := 0
	err = stmt.QueryRow(
		recordId,
		r.data.ChatID,
		r.data.MessageHash,
		r.data.MessageID,
		r.data.ReceiverKeyUID,
		r.data.SentAt,
		r.data.Topic,
		r.data.MessageType,
		r.data.MessageSize,
		r.data.PubsubTopic,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func QueryReceivedMessagesBetween(db *sql.DB, startsAt time.Time, endsAt time.Time) ([]*types.ReceivedMessage, error) {
	rows, err := db.Query(`
	SELECT rm.id, rm.chatId, rm.messageHash, rm.messageId, rm.receiverKeyUID, rm.sentAt, rm.topic, rm.messageType, rm.messageSize, rm.pubSubTopic,
		   cf.nodeName, cf.peerId, cf.statusVersion, cf.deviceType
	FROM receivedMessages rm
	LEFT JOIN telemetryRecord cf ON rm.recordId = cf.id
	WHERE rm.sentAt BETWEEN $1 AND $2`, startsAt.Unix(), endsAt.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivedMessages []*types.ReceivedMessage
	for rows.Next() {
		var receivedMessage types.ReceivedMessage
		err = rows.Scan(
			&receivedMessage.ID,
			&receivedMessage.ChatID,
			&receivedMessage.MessageHash,
			&receivedMessage.MessageID,
			&receivedMessage.ReceiverKeyUID,
			&receivedMessage.SentAt,
			&receivedMessage.Topic,
			&receivedMessage.MessageType,
			&receivedMessage.MessageSize,
			&receivedMessage.PubsubTopic,
			&receivedMessage.NodeName,
			&receivedMessage.PeerID,
			&receivedMessage.StatusVersion,
			&receivedMessage.DeviceType,
		)
		if err != nil {
			return nil, err
		}
		receivedMessages = append(receivedMessages, &receivedMessage)
	}
	return receivedMessages, nil
}

func DidReceivedMessageBeforeAndAfterInChat(db *sql.DB, receiverPublicKey string, before, after time.Time, chatId string) (bool, error) {
	var afterCount int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM receivedMessages rm
		JOIN telemetryRecord cf ON rm.recordId = cf.id
		WHERE rm.receiverKeyUID = $1 AND cf.createdAt > $2 AND rm.chatId = $3`,
		receiverPublicKey,
		after.Unix(),
		chatId,
	).Scan(&afterCount)
	if err != nil {
		return false, err
	}

	var beforeCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM receivedMessages rm
		JOIN telemetryRecord cf ON rm.recordId = cf.id
		WHERE rm.receiverKeyUID = $1 AND cf.createdAt < $2 AND rm.chatId = $3`,
		receiverPublicKey,
		before.Unix(),
		chatId,
	).Scan(&beforeCount)
	if err != nil {
		return false, err
	}

	return afterCount > 0 && beforeCount > 0, nil
}

func (r *ReceivedMessageAggregated) Put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO receivedMessageAggregated (chatId, durationInSeconds, value, runAt) VALUES ($1, $2, $3, $4) RETURNING id;")
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(r.ChatID, r.DurationInSeconds, r.Value, r.RunAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}
