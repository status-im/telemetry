package metrics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

type ProtocolStats struct {
	types.ProtocolStats
}

func (r *ProtocolStats) insertRate(db *sql.DB, protocolName string, metric types.Metric) error {
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

func (r *ProtocolStats) Put(db *sql.DB) error {
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

func (r *ProtocolStats) Process(db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) (err error) {
	if err := json.Unmarshal(*data.TelemetryData, &r); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error decoding protocol stats: %v", err))
		return err
	}

	if err := r.Put(db); err != nil {
		errs.Append(data.Id, fmt.Sprintf("Error saving protocol stats: %v", err))
		return err
	}

	return nil
}

func (r *ProtocolStats) Clean(db *sql.DB, before int64) (int64, error) {
	result, err := db.Exec("DELETE FROM protocolStatsRate WHERE createdAt < $1", before)
	if err != nil {
		return 0, err
	}

	rows1, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	result2, err := db.Exec("DELETE FROM protocolStatsTotals WHERE createdAt < $1", time.Unix(before, 0))
	if err != nil {
		return 0, err
	}

	rows2, err := result2.RowsAffected()

	return rows1 + rows2, err
}
