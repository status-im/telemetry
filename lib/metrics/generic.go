package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
)

// GenericMetric is a generic struct that can handle any metric type
type GenericMetric[T any] struct {
	types.TelemetryRecord
	Data T
}

// MetricProcessor is an interface for processing metrics
type MetricProcessor interface {
	Process(context.Context, *sql.DB, *common.MetricErrors, *types.TelemetryRequest) error
	Clean(*sql.DB, int64) (int64, error)
}

// NewMetricProcessor creates a new MetricProcessor for the given metric type
func NewMetricProcessor[T types.TelemetryRecord]() MetricProcessor {
	return &GenericMetric[T]{
		Data: *new(T),
	}
}

// Process implements the MetricProcessor interface
func (g *GenericMetric[T]) Process(ctx context.Context, db *sql.DB, errs *common.MetricErrors, data *types.TelemetryRequest) error {
	// Unmarshal the TelemetryRecord fields
	if err := json.Unmarshal(*data.TelemetryData, &g.TelemetryRecord); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding TelemetryRecord: %v", err))
		return err
	}

	// Unmarshal the Data field
	if err := json.Unmarshal(*data.TelemetryData, &g.Data); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error decoding %T: %v", g.Data, err))
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recordId, err := InsertTelemetryRecord(tx, &g.TelemetryRecord)
	if err != nil {
		return err
	}

	columns, values := getColumnsAndValues(g.Data)
	placeholders := make([]string, len(columns)+1)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	tableName := getTableName(g.Data)
	query := fmt.Sprintf(`
		INSERT INTO %s (recordId, %s)
		VALUES (%s)
		RETURNING id;
	`, tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	args := []interface{}{recordId}
	args = append(args, values...)

	result := tx.QueryRowContext(ctx, query, args...)
	if result.Err() != nil {
		errs.Append(data.ID, fmt.Sprintf("Error saving %T: %v", g.Data, result.Err()))
		return result.Err()
	}

	if err := tx.Commit(); err != nil {
		errs.Append(data.ID, fmt.Sprintf("Error committing transaction: %v", err))
		return err
	}

	return nil
}

// Clean implements the MetricProcessor interface
func (g *GenericMetric[T]) Clean(db *sql.DB, before int64) (int64, error) {
	tableName := getTableName(g.Data)
	return common.Cleanup(db, tableName, before)
}

// Helper functions

func getColumnsAndValues(v interface{}) ([]string, []interface{}) {
	var columns []string
	var values []interface{}
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			columnName := strings.Split(tag, ",")[0]
			columns = append(columns, columnName)
			values = append(values, val.Field(i).Interface())
		}
	}
	return columns, values
}

func getTableName(v interface{}) string {
	t := reflect.TypeOf(v)
	return strings.ToLower(t.Name())
}
