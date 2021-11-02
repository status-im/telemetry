package telemetry

import (
	"database/sql"
	"time"
)

type ReceivedMessage struct {
	ID             int    `json:"id"`
	ChatId         string `json:"chatId"`
	MessageHash    string `json:"messageHash"`
	ReceiverKeyUID string `json:"receiverKeyUID"`
	SentAt         int64  `json:"sentAt"`
	Topic          string `json:"topic"`
	CreatedAt      int64  `json:"createdAt"`
}

func (s *ReceivedMessage) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO receivedMessages (chatId, messageHash, receiverKeyUID, sentAt, topic, createdAt) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	s.CreatedAt = time.Now().Unix()
	res, err := stmt.Exec(s.ChatId, s.MessageHash, s.ReceiverKeyUID, s.SentAt, s.Topic, s.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)
	return nil
}
