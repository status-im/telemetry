package metrics

import (
	"github.com/status-im/telemetry/pkg/types"
)

type MissedMessages struct {
	GenericMetric[types.MissedMessages]
}

type MissedRelevantMessages struct {
	GenericMetric[types.MissedRelevantMessages]
}

type MessageDeliveryConfirmed struct {
	GenericMetric[types.MessageDeliveryConfirmed]
}

var (
	MissedMessagesProcessor           = &MissedMessages{}
	MissedRelevantMessagesProcessor   = &MissedRelevantMessages{}
	MessageDeliveryConfirmedProcessor = &MessageDeliveryConfirmed{}
)
