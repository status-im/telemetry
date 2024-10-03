package metrics

import (
	"github.com/status-im/telemetry/pkg/types"
)

type MissedMessage struct {
	GenericMetric[types.MissedMessage]
}

type MissedRelevantMessage struct {
	GenericMetric[types.MissedRelevantMessage]
}

type MessageDeliveryConfirmed struct {
	GenericMetric[types.MessageDeliveryConfirmed]
}

var (
	MissedMessageProcessor            = &MissedMessage{}
	MissedRelevantMessageProcessor    = &MissedRelevantMessage{}
	MessageDeliveryConfirmedProcessor = &MessageDeliveryConfirmed{}
)
