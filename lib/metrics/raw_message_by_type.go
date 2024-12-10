package metrics

import "github.com/status-im/telemetry/pkg/types"

type RawMessageByTypeMetric struct {
	GenericMetric[types.RawMessageByType]
}
