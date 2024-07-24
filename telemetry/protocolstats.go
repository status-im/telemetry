package telemetry

import (
	"database/sql"
	"time"
)

type Metric struct {
	TotalIn  int64   `json:"totalIn"`
	TotalOut int64   `json:"totalOut"`
	RateIn   float64 `json:"rateIn"`
	RateOut  float64 `json:"rateOut"`
}

type ProtocolStats struct {
	PeerID          string `json:"hostID"`
	Relay           Metric `json:"relay"`
	Store           Metric `json:"store"`
	FilterPush      Metric `json:"filter-push"`
	FilterSubscribe Metric `json:"filter-subscribe"`
	Lightpush       Metric `json:"lightpush"`
}

func (r *ProtocolStats) insertRate(db *sql.DB, protocolName string, metric Metric) error {
	stmt, err := db.Prepare("INSERT INTO protocolStatsRate (peerID, protocolName, rateIn, rateOut, createdAt) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.PeerID, protocolName, metric.RateIn, metric.RateOut, time.Now().Unix())
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("INSERT INTO protocolStatsTotals (peerID, protocolName, totalIn, totalOut, createdAt) VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT protocolStatsTotals_unique DO UPDATE SET totalIn = $3, totalOut = $4;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.PeerID, protocolName, metric.TotalIn, metric.TotalOut, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	return nil
}

func (r *ProtocolStats) put(db *sql.DB) error {
	err := r.insertRate(db, "relay", r.Relay)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "store", r.Store)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "filter-push", r.FilterPush)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "filter-subscribe", r.FilterSubscribe)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "lightpush", r.Lightpush)
	if err != nil {
		return err
	}

	return nil
}
