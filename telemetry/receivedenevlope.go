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
	PeerID          string `json:"peerId"`
	NodeName        string `json:"nodeName"`
	ProcessingError string `json:"processingError"`
	StatusVersion   string `json:"statusVersion"`
}

func (r *ReceivedEnvelope) put(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO receivedEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, receiverKeyUID, nodeName, processingError, statusVersion)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
							ON CONFLICT ON CONSTRAINT receivedEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := 0
	err = stmt.QueryRow(r.MessageHash, r.SentAt, r.CreatedAt, r.PubsubTopic, r.Topic, r.ReceiverKeyUID, r.NodeName, r.ProcessingError, r.StatusVersion).Scan(&lastInsertId)
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

type SentEnvelope struct {
	ID              int    `json:"id"`
	MessageHash     string `json:"messageHash"`
	SentAt          int64  `json:"sentAt"`
	CreatedAt       int64  `json:"createdAt"`
	PubsubTopic     string `json:"pubsubTopic"`
	Topic           string `json:"topic"`
	SenderKeyUID    string `json:"senderKeyUID"`
	PeerID          string `json:"peerId"`
	NodeName        string `json:"nodeName"`
	ProcessingError string `json:"processingError"`
	PublishMethod   string `json:"publishMethod"`
	StatusVersion   string `json:"statusVersion"`
}

func (r *SentEnvelope) put(db *sql.DB) error {
	r.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO sentEnvelopes (messageHash, sentAt, createdAt, pubsubTopic,
							topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
							ON CONFLICT ON CONSTRAINT sentEnvelopes_unique DO NOTHING
							RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(r.MessageHash, r.SentAt, r.CreatedAt, r.PubsubTopic, r.Topic, r.SenderKeyUID, r.PeerID, r.NodeName, r.PublishMethod, r.StatusVersion).Scan(&lastInsertId)
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

type ErrorSendingEnvelope struct {
	CreatedAt    int64        `json:"createdAt"`
	Error        string       `json:"error"`
	SentEnvelope SentEnvelope `json:"sentEnvelope"`
}

func (e *ErrorSendingEnvelope) put(db *sql.DB) error {
	e.CreatedAt = time.Now().Unix()
	stmt, err := db.Prepare(`INSERT INTO errorSendingEnvelope (messageHash, sentAt, createdAt, pubsubTopic,
		topic, senderKeyUID, peerId, nodeName, publishMethod, statusVersion, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT ON CONSTRAINT errorSendingEnvelope_unique DO NOTHING
		RETURNING id;`)
	if err != nil {
		return err
	}

	lastInsertId := int64(0)
	err = stmt.QueryRow(e.SentEnvelope.MessageHash, e.SentEnvelope.SentAt, e.CreatedAt, e.SentEnvelope.PubsubTopic, e.SentEnvelope.Topic, e.SentEnvelope.SenderKeyUID, e.SentEnvelope.PeerID, e.SentEnvelope.NodeName, e.SentEnvelope.PublishMethod, e.SentEnvelope.StatusVersion, e.Error).Scan(&lastInsertId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	defer stmt.Close()
	e.SentEnvelope.ID = int(lastInsertId)

	return nil
}
