package metrics

import (
	"database/sql"
	"encoding/json"
	"time"
)

type WakuTelemetryType string

const (
	LightPushFilter WakuTelemetryType = "LightPushFilter"
	Generic         WakuTelemetryType = "Generic"
)

type WakuTelemetryRequest struct {
	Id            int               `json:"id"`
	TelemetryType WakuTelemetryType `json:"telemetryType"`
	TelemetryData *json.RawMessage  `json:"telemetryData"`
}

type TelemetryPushFilter struct {
	ID            int    `json:"id"`
	Protocol      string `json:"protocol"`
	Ephemeral     bool   `json:"ephemeral"`
	Timestamp     int64  `json:"timestamp"`
	SeenTimestamp int64  `json:"seenTimestamp"`
	CreatedAt     int64  `json:"createdAt"`
	ContentTopic  string `json:"contentTopic"`
	PubsubTopic   string `json:"pubsubTopic"`
	PeerID        string `json:"peerId"`
	MessageHash   string `json:"messageHash"`
	ErrorMessage  string `json:"errorMessage"`
	ExtraData     string `json:"extraData"`
}

func (r *TelemetryPushFilter) Put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO wakuPushFilter (protocol, ephemeral, timestamp, seenTimestamp, contentTopic, pubsubTopic, peerId, messageHash, errorMessage, extraData, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(
		r.Protocol,
		r.Ephemeral,
		r.Timestamp,
		r.SeenTimestamp,
		r.ContentTopic,
		r.PubsubTopic,
		r.PeerID,
		r.MessageHash,
		r.ErrorMessage,
		r.ExtraData,
		r.CreatedAt,
	).Scan(&lastInsertId)

	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}

type TelemetryGeneric struct {
	ID           int    `json:"id"`
	PeerID       string `json:"peerId"`
	MetricType   string `json:"metricType"`
	ContentTopic string `json:"contentTopic"`
	PubsubTopic  string `json:"pubsubTopic"`
	GenericData  string `json:"genericData"`
	ErrorMessage string `json:"errorMessage"`
	Timestamp    int64  `json:"timestamp"`
	CreatedAt    int64  `json:"createdAt"`
}

func (r *TelemetryGeneric) Put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO wakuGeneric (peerId, metricType, contentTopic, pubsubTopic, genericData, errorMessage, timestamp, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.PeerID, r.MetricType, r.ContentTopic, r.PubsubTopic, r.GenericData, r.ErrorMessage, r.Timestamp, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}
