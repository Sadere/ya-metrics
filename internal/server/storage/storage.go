package storage

import (
	"github.com/Sadere/ya-metrics/internal/common"
)

// Интерфейс для хранения данных о метриках
type MetricRepository interface {
	Get(common.MetricType, string) (common.Metric, error)
	Set(string, common.Metric) error

	GetData() map[string]common.Metric
}
