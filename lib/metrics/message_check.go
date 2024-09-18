package metrics

import (
	"github.com/status-im/telemetry/pkg/types"
)

type MessageCheckSuccess struct {
	GenericMetric[types.MessageCheckSuccess]
}

type MessageCheckFailure struct {
	GenericMetric[types.MessageCheckFailure]
}

var (
	MessageCheckSuccessProcessor = &MessageCheckSuccess{}
	MessageCheckFailureProcessor = &MessageCheckFailure{}
)
