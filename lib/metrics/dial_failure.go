package metrics

import (
	"github.com/status-im/telemetry/pkg/types"
)

type DialFailure struct {
	GenericMetric[types.DialFailure]
}

var (
	DialFailureProcessor = &DialFailure{}
)
