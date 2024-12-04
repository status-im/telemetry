package metrics

import "github.com/status-im/telemetry/pkg/types"

type SentMessageTotalMetric struct {
	GenericMetric[types.SentMessageTotal]
}
