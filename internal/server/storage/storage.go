package storage

import (
	"github.com/Sadere/ya-metrics/internal/common"
)

// Интерфейс для хранения данных о метриках
type MetricRepository interface {
	Get(common.MetricType, string) (common.Metrics, error)
	Set(string, common.Metrics) error

	GetData() (map[string]common.Metrics, error)
	SetData(map[string]common.Metrics) error
}
