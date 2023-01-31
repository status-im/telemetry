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
	HostID string `json:"hostID"`
	Relay  Metric `json:"relay"`
	Store  Metric `json:"store"`
}

func (r *ProtocolStats) insertRate(db *sql.DB, protocolName string, metric Metric) error {
	stmt, err := db.Prepare("INSERT INTO protocolStatsRate (hostId, protocolName, rateIn, rateOut, createdAt) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.HostID, protocolName, metric.RateIn, metric.RateOut, time.Now().Unix())
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("INSERT INTO protocolStatsTotals (hostId, protocolName, totalIn, totalOut, createdAt) VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT protocolStatsTotals_unique DO UPDATE SET totalIn = $3, totalOut = $4;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.HostID, protocolName, metric.TotalIn, metric.TotalOut, time.Now().Format("2006-01-02"))
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

	return nil
}
