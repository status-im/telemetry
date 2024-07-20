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
	PeerCount     int    `json:"peerCount"`
	StatusVersion string `json:"statusVersion"`
}

func (r *PeerCount) put(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO peerCount (timestamp, nodeName, nodeKeyUid, peerCount, statusVersion, createdAt) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r.CreatedAt = time.Now().Unix()
	lastInsertId := 0
	err = stmt.QueryRow(r.Timestamp, r.NodeName, r.NodeKeyUid, r.PeerCount, r.StatusVersion, r.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	r.ID = lastInsertId

	return nil
}