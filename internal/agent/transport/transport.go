package transport

import "github.com/Sadere/ya-metrics/internal/common"

type MetricTransport interface {
	SendMetrics([]common.Metrics) error
}
