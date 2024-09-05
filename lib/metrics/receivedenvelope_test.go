package metrics

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEnvelopesUpdate(t *testing.T) {
	db := common.NewMock()
	defer common.DropTables(db)

	var errs common.MetricErrors

	firstEnvelopeData := types.ReceivedEnvelope{
		TelemetryRecord: types.TelemetryRecord{
			NodeName: "status",
		},
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
		PubsubTopic:    "1",
	}

	data, err := json.Marshal(firstEnvelopeData)
	require.NoError(t, err)

	telemetryRequest1 := types.TelemetryRequest{
		ID:            0,
		TelemetryType: types.ReceivedEnvelopeMetric,
		TelemetryData: (*json.RawMessage)(&data),
	}

	var firstEnvelope ReceivedEnvelope
	err = firstEnvelope.Process(context.Background(), db, &errs, &telemetryRequest1)
	require.NoError(t, err)

	envelopeToUpdateData := types.ReceivedEnvelope{
		TelemetryRecord: types.TelemetryRecord{
			NodeName: "status",
		},
		MessageHash:     "1",
		ReceiverKeyUID:  "1",
		SentAt:          time.Now().Unix(),
		Topic:           "1",
		PubsubTopic:     "1",
		ProcessingError: "MyError",
	}
	envelopeToUpdate := ReceivedEnvelope{
		envelopeToUpdateData,
	}

	err = envelopeToUpdate.UpdateProcessingError(db)
	require.NoError(t, err)

	rows, err := db.Query("SELECT processingerror FROM receivedEnvelopes WHERE messagehash = '1';")
	require.NoError(t, err)
	defer rows.Close()
	rows.Next()
	var procError = ""
	err = rows.Scan(&procError)
	require.NoError(t, err)
	require.Equal(t, "MyError", procError)
}
