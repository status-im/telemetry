package common

import (
	"database/sql"

	"github.com/status-im/telemetry/pkg/types"
)

type MetricProcessor interface {
	Process(db *sql.DB, errs *MetricErrors, data *types.TelemetryRequest) error
	Clean(db *sql.DB, before int64) (int64, error)
}
