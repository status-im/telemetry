package telemetry

import (
	"database/sql"
	"encoding/json"
	"time"
)

type WakuTelemetryType string

const (
	LightPushFilter WakuTelemetryType = "LightPushFilter"
)

type WakuTelemetryRequest struct {
	Id            int               `json:"id"`
	TelemetryType WakuTelemetryType `json:"telemetryType"`
	TelemetryData *json.RawMessage  `json:"telemetryData"`
}

type TelemetryPushFilter struct {
	ID             int    `json:"id"`
	WalletAddress  string `json:"walletAddress"`
	PeerIDSender   string `json:"peerIdSender"`
	PeerIDReporter string `json:"peerIdReporter"`
	SequenceHash   string `json:"sequenceHash"`
	SequenceTotal  uint64 `json:"sequenceTotal"`
	SequenceIndex  uint64 `json:"sequenceIndex"`
	ContentTopic   string `json:"contentTopic"`
	PubsubTopic    string `json:"pubsubTopic"`
	Timestamp      int64  `json:"timestamp"`
	CreatedAt      int64  `json:"createdAt"`
}

func (r *TelemetryPushFilter) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO wakuPushFilter (peerIdSender, peerIdReporter, sequenceHash, sequenceTotal, sequenceIndex, contentTopic, pubsubTopic, timestamp, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.PeerIDSender, r.PeerIDReporter, r.SequenceHash, r.SequenceTotal, r.SequenceIndex, r.ContentTopic, r.PubsubTopic, r.Timestamp, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}
