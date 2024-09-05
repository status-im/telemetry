package telemetry

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/lib/metrics"
	"github.com/status-im/telemetry/pkg/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func updateCreatedAt(db *sql.DB, m *metrics.ReceivedMessage, createdAt int64) error {
	_, err := db.Exec(`
		UPDATE telemetryRecord
		SET createdAt = $1
		WHERE id = (
			SELECT recordId
			FROM receivedMessages
			WHERE id = $2
		)
	`, createdAt, m.ID)
	return err
}

func queryAggregatedMessage(db *sql.DB) ([]*metrics.ReceivedMessageAggregated, error) {
	rows, err := db.Query("SELECT * FROM receivedMessageAggregated")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivedMessageAggregateds []*metrics.ReceivedMessageAggregated
	for rows.Next() {
		var receivedMessageAggregated metrics.ReceivedMessageAggregated
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
	db := common.NewMock()
	defer common.DropTables(db)

	mData := types.ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}

	m := &metrics.ReceivedMessage{ReceivedMessage: mData}
	ctx := context.Background()
	err := m.Put(ctx, db)
	require.NoError(t, err)

	oneHourAndHalf := time.Hour + time.Minute*30
	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)

	twoHourAndHalf := 5*time.Hour + time.Minute*30
	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "3",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)
	err = updateCreatedAt(db, m, m.SentAt)
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
	db := common.NewMock()
	defer common.DropTables(db)

	mData := types.ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}
	m := &metrics.ReceivedMessage{ReceivedMessage: mData}
	ctx := context.Background()
	err := m.Put(ctx, db)
	require.NoError(t, err)

	oneHourAndHalf := time.Hour + time.Minute*30
	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)

	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "3",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)

	twoHourAndHalf := 5*time.Hour + time.Minute*30
	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "4",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)
	err = updateCreatedAt(db, m, m.SentAt)
	require.NoError(t, err)

	mData = types.ReceivedMessage{
		ChatID:         "1",
		MessageHash:    "1",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)

	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "2",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Add(-oneHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)

	mData = types.ReceivedMessage{
		ChatID:         "3",
		MessageHash:    "4",
		ReceiverKeyUID: "2",
		SentAt:         time.Now().Add(-twoHourAndHalf).Unix(),
		Topic:          "1",
	}
	m = &metrics.ReceivedMessage{ReceivedMessage: mData}
	err = m.Put(ctx, db)
	require.NoError(t, err)
	err = updateCreatedAt(db, m, m.SentAt)
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
