package telemetry

import (
	"database/sql"
	"log"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func NewMock() *sql.DB {
	db, err := sql.Open("postgres", "postgres://telemetry:newPassword@127.0.0.1:5432/telemetrydb?sslmode=disable")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	err = createTables(db)

	if err != nil {
		log.Fatalf("an error '%s' was not expected when migrating the db", err)
	}

	return db
}

func dropTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS receivedMessages")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS receivedMessageAggregated")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS receivedEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS protocolStatsRate")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS protocolStatsTotals")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS schema_migrations")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS receivedEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS receivedMessageAggregated_runAt")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS protocolStatsRate_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS protocolStatsTotals_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	db.Close()
}

func updateCreatedAt(db *sql.DB, m *ReceivedMessage) error {
	_, err := db.Exec("UPDATE receivedMessages SET createdAt = $1 WHERE id = $2", m.CreatedAt, m.ID)
	return err
}

func queryAggregatedMessage(db *sql.DB) ([]*ReceivedMessageAggregated, error) {
	rows, err := db.Query("SELECT * FROM receivedMessageAggregated")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivedMessageAggregateds []*ReceivedMessageAggregated
	for rows.Next() {
		var receivedMessageAggregated ReceivedMessageAggregated
		err = rows.Scan(
			&receivedMessageAggregated.ID,
			&receivedMessageAggregated.DurationInSeconds,
			&receivedMessageAggregated.ChatID,
			&receivedMessageAggregated.Value,
			&receivedMessageAggregated.RunAt,
		)
		if err != nil {
			return nil, err
		}
		receivedMessageAggregateds = append(receivedMessageAggregateds, &receivedMessageAggregated)
	}
	return receivedMessageAggregateds, nil
}

func TestRunAggregatorSimple(t *testing.T) {
	db := NewMock()
	defer dropTables(db)

	m := &ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}
	err := m.put(db)
	require.NoError(t, err)

	oneHourAndHalf := time.Hour + time.Minute*30
	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)

	twoHourAndHalf := 5*time.Hour + time.Minute*30
	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "3",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)
	m.CreatedAt = m.SentAt
	err = updateCreatedAt(db, m)
	require.NoError(t, err)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	agg, err := NewAggregator(db, logger)
	require.NoError(t, err)

	agg.Run(time.Hour)

	res, err := queryAggregatedMessage(db)
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, "3", res[0].ChatID)
	require.Equal(t, 1.0, res[0].Value)
	require.Equal(t, "", res[1].ChatID)
	require.Equal(t, 1.0, res[1].Value)
}

func TestRunAggregatorSimpleWithMessageMissing(t *testing.T) {
	db := NewMock()
	defer dropTables(db)

	m := &ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}
	err := m.put(db)
	require.NoError(t, err)

	oneHourAndHalf := time.Hour + time.Minute*30
	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)

	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "3",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)

	twoHourAndHalf := 5*time.Hour + time.Minute*30
	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "4",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)
	m.CreatedAt = m.SentAt
	err = updateCreatedAt(db, m)
	require.NoError(t, err)

	m = &ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)

	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)

	m = &ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "4",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	err = m.put(db)
	require.NoError(t, err)
	m.CreatedAt = m.SentAt
	err = updateCreatedAt(db, m)
	require.NoError(t, err)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	agg, err := NewAggregator(db, logger)
	require.NoError(t, err)

	agg.Run(time.Hour)

	res, err := queryAggregatedMessage(db)
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, "3", res[0].ChatID)
	require.Equal(t, 0.67, math.Round(res[0].Value*100)/100)
	require.Equal(t, "", res[1].ChatID)
	require.Equal(t, 0.67, math.Round(res[1].Value*100)/100)
}
