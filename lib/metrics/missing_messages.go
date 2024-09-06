package metrics

import (
	"github.com/status-im/telemetry/pkg/types"
)

type MissingMessage struct {
	GenericMetric[types.MissingMessage]
}

type MissingRelevantMessage struct {
	GenericMetric[types.MissingRelevantMessage]
}

var (
	MissingMessageProcessor         = &MissingMessage{}
	MissingRelevantMessageProcessor = &MissingRelevantMessage{}
)