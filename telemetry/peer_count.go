package telemetry

import (
	"database/sql"
	"time"
)

type PeerCount struct {
	ID            int    `json:"id"`
	CreatedAt     int64  `json:"createdAt"`
	Timestamp     int64  `json:"timestamp"`
	NodeName      string `json:"nodeName"`
	NodeKeyUid    string `json:"nodeKeyUid"`
	PeerID        string `json:"peerId"`
	PeerCount     int    `json:"peerCount"`
	StatusVersion string `json:"statusVersion"`
}

func (r *PeerCount) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO peerCount (timestamp, nodeName, nodeKeyUid, peerId, peerCount, statusVersion, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.Timestamp, r.NodeName, r.NodeKeyUid, r.PeerID, r.PeerCount, r.StatusVersion, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}

type PeerConnFailure struct {
	ID            int    `json:"id"`
	CreatedAt     int64  `json:"createdAt"`
	Timestamp     int64  `json:"timestamp"`
	NodeName      string `json:"nodeName"`
	NodeKeyUid    string `json:"nodeKeyUid"`
	PeerId        string `json:"peerId"`
	StatusVersion string `json:"statusVersion"`
	FailedPeerId  string `json:"failedPeerId"`
	FailureCount  int    `json:"failureCount"`
}

func (r *PeerConnFailure) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO peerConnFailure (timestamp, nodeName, nodeKeyUid, peerId, failedPeerId, failureCount, statusVersion, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.Timestamp, r.NodeName, r.NodeKeyUid, r.PeerId, r.FailedPeerId, r.FailureCount, r.StatusVersion, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}
