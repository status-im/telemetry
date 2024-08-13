package telemetry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/status-im/dev-telemetry/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEnvelopesUpdate(t *testing.T) {
	db := NewMock()
	defer dropTables(db)

	var errs MetricErrors

	firstEnvelopeData := types.ReceivedEnvelope{
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		NodeName:       "status",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
		PubsubTopic:    "1",
	}

	data, err := json.Marshal(firstEnvelopeData)
	require.NoError(t, err)

	telemetryRequest1 := types.TelemetryRequest{
		Id:            0,
		TelemetryType: types.ReceivedEnvelopeMetric,
		TelemetryData: (*json.RawMessage)(&data),
	}

	var firstEnvelope ReceivedEnvelope
	err = firstEnvelope.process(db, &errs, &telemetryRequest1)
	require.NoError(t, err)

	envelopeToUpdateData := types.ReceivedEnvelope{
		MessageHash:     "1",
		ReceiverKeyUID:  "1",
		NodeName:        "status",
		SentAt:          time.Now().Unix(),
		Topic:           "1",
		PubsubTopic:     "1",
		ProcessingError: "MyError",
	}
	envelopeToUpdate := ReceivedEnvelope{
		data: envelopeToUpdateData,
	}

	err = envelopeToUpdate.updateProcessingError(db)
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
