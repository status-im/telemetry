package telemetry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
	data types.ReceivedMessage
}

func (r *ReceivedMessage) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding received message failure: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving received messages: %v", err))
		return err
	}
	return nil
}

func (r *ReceivedMessage) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO receivedMessages (chatId, messageHash, messageId, receiverKeyUID, peerId, nodeName, sentAt, topic, messageType, messageSize, createdAt, pubSubTopic, statusVersion, deviceType) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id;")
	if err != nil {
		return err
	}

	r.data.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(
		r.data.ChatID,
		r.data.MessageHash,
		r.data.MessageID,
		r.data.ReceiverKeyUID,
		r.data.PeerID,
		r.data.NodeName,
		r.data.SentAt,
		r.data.Topic,
		r.data.MessageType,
		r.data.MessageSize,
		r.data.CreatedAt,
		r.data.PubsubTopic,
		r.data.StatusVersion,
		r.data.DeviceType,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.data.ID = lastInsertId

	return nil
}

func queryReceivedMessagesBetween(db *sql.DB, startsAt time.Time, endsAt time.Time) ([]*types.ReceivedMessage, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, chatId, messageHash, messageId, receiverKeyUID, peerId, nodeName, sentAt, topic, messageType, messageSize, createdAt, pubSubTopic FROM receivedMessages WHERE sentAt BETWEEN %d and %d", startsAt.Unix(), endsAt.Unix()))
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
			&receivedMessage.PeerID,
			&receivedMessage.NodeName,
			&receivedMessage.SentAt,
			&receivedMessage.Topic,
			&receivedMessage.MessageType,
			&receivedMessage.MessageSize,
			&receivedMessage.CreatedAt,
			&receivedMessage.PubsubTopic,
		)
		if err != nil {
			return nil, err
		}
		receivedMessages = append(receivedMessages, &receivedMessage)
	}
	return receivedMessages, nil
}

func didReceivedMessageBeforeAndAfterInChat(db *sql.DB, receiverPublicKey string, before, after time.Time, chatId string) (bool, error) {
	var afterCount int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM receivedMessages WHERE receiverKeyUID = $1 AND createdAt > $2 AND chatId = $3",
		receiverPublicKey,
		after.Unix(),
		chatId,
	).Scan(&afterCount)
	if err != nil {
		return false, err
	}

	var beforeCount int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM receivedMessages WHERE receiverKeyUID = $1 AND createdAt < $2 AND chatId = $3",
		receiverPublicKey,
		before.Unix(),
		chatId,
	).Scan(&beforeCount)
	if err != nil {
		return false, err
	}

	return afterCount > 0 && beforeCount > 0, nil
}

func (r *ReceivedMessageAggregated) put(db *sql.DB) error {
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
