package telemetry

import (
	"database/sql"
	"errors"
	"time"
)

type ReceivedEnvelope struct {
	ID              int    `json:"id"`
	MessageHash     string `json:"messageHash"`
	SentAt          int64  `json:"sentAt"`
	CreatedAt       int64  `json:"createdAt"`
	PubsubTopic     string `json:"pubsubTopic"`
	Topic           string `json:"topic"`
	ReceiverKeyUID  string `json:"receiverKeyUID"`
	NodeName        string `json:"nodeName"`
	ProcessingError string `json:"processingError"`
}

func (r *ReceivedEnvelope) put(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO receivedEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, receiverKeyUID, nodeName, processingError)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(r.MessageHash, r.SentAt, r.CreatedAt, r.PubsubTopic, r.Topic, r.ReceiverKeyUID, r.NodeName, r.ProcessingError).Scan(&lastInsertId)
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

func (r *ReceivedEnvelope) updateProcessingError(db *sql.DB) error {
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
