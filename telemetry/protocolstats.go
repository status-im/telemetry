package telemetry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/telemetry/pkg/types"
)

type ProtocolStats struct {
	data types.ProtocolStats
}

func (r *ProtocolStats) insertRate(db *sql.DB, protocolName string, metric types.Metric) error {
	stmt, err := db.Prepare("INSERT INTO protocolStatsRate (peerID, protocolName, rateIn, rateOut, createdAt) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.data.PeerID, protocolName, metric.RateIn, metric.RateOut, time.Now().Unix())
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("INSERT INTO protocolStatsTotals (peerID, protocolName, totalIn, totalOut, createdAt) VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT protocolStatsTotals_unique DO UPDATE SET totalIn = $3, totalOut = $4;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.data.PeerID, protocolName, metric.TotalIn, metric.TotalOut, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	return nil
}

func (r *ProtocolStats) put(db *sql.DB) error {
	err := r.insertRate(db, "relay", r.data.Relay)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "store", r.data.Store)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "filter-push", r.data.FilterPush)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "filter-subscribe", r.data.FilterSubscribe)
	if err != nil {
		return err
	}

	err = r.insertRate(db, "lightpush", r.data.Lightpush)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProtocolStats) process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) (err error) {
	if err := json.Unmarshal(*data.TelemetryData, &r.data); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding protocol stats: %v", err))
		return err
	}

	if err := r.put(db); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving protocol stats: %v", err))
		return err
	}

	return nil
}
