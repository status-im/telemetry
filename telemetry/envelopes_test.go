package telemetry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnvelopesUpdate(t *testing.T) {
	db := NewMock()
	defer dropTables(db)

	firstEnvelope := &ReceivedEnvelope{
		MessageHash:    "1",
		ReceiverKeyUID: "1",
		NodeName:       "status",
		SentAt:         time.Now().Unix(),
		Topic:          "1",
		PubsubTopic:    "1",
	}
	err := firstEnvelope.put(db)
	require.NoError(t, err)

	envelopeToUpdate := &ReceivedEnvelope{
		MessageHash:     "1",
		ReceiverKeyUID:  "1",
		NodeName:        "status",
		SentAt:          time.Now().Unix(),
		Topic:           "1",
		PubsubTopic:     "1",
		ProcessingError: "MyError",
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
