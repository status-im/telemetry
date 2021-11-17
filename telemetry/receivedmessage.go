package telemetry

import (
	"database/sql"
	"fmt"
	"time"
)

type ReceivedMessageAggregated struct {
	ID                int
	ChatID            string
	DurationInSeconds int64
	Value             float64
	RunAt             int64
}

type ReceivedMessage struct {
	ID             int    `json:"id"`
	ChatID         string `json:"chatId"`
	MessageHash    string `json:"messageHash"`
	MessageID      string `json:"messageId"`
	ReceiverKeyUID string `json:"receiverKeyUID"`
	NodeName       string `json:"nodeName"`
	SentAt         int64  `json:"sentAt"`
	Topic          string `json:"topic"`
	CreatedAt      int64  `json:"createdAt"`
}

func queryReceivedMessagesBetween(db *sql.DB, startsAt time.Time, endsAt time.Time) ([]*ReceivedMessage, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, chatId, messageHash, messageId, receiverKeyUID, nodeName, sentAt, topic, createdAt FROM receivedMessages WHERE sentAt BETWEEN %d and %d", startsAt.Unix(), endsAt.Unix()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivedMessages []*ReceivedMessage
	for rows.Next() {
		var receivedMessage ReceivedMessage
		err = rows.Scan(
			&receivedMessage.ID,
			&receivedMessage.ChatID,
			&receivedMessage.MessageHash,
			&receivedMessage.MessageID,
			&receivedMessage.ReceiverKeyUID,
			&receivedMessage.NodeName,
			&receivedMessage.SentAt,
			&receivedMessage.Topic,
			&receivedMessage.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		receivedMessages = append(receivedMessages, &receivedMessage)
	}
	return receivedMessages, nil
}

func didReceivedMessageBeforeAndAfter(db *sql.DB, receiverPublicKey string, before, after time.Time) (bool, error) {
	var afterCount int
	err := db.QueryRow("SELECT COUNT(*) FROM receivedMessages WHERE receiverKeyUID = $1 AND createdAt > $2", receiverPublicKey, after.Unix()).Scan(&afterCount)
	if err != nil {
		return false, err
	}

	var beforeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM receivedMessages WHERE receiverKeyUID = $1 AND createdAt < $2", receiverPublicKey, before.Unix()).Scan(&beforeCount)
	if err != nil {
		return false, err
	}

	return afterCount > 0 && beforeCount > 0, nil
}

func (r *ReceivedMessage) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO receivedMessages (chatId, messageHash, messageId, receiverKeyUID, nodeName, sentAt, topic, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
	if err != nil {
		return err
	}

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.ChatID, r.MessageHash, r.MessageID, r.ReceiverKeyUID, r.NodeName, r.SentAt, r.Topic, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
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
