package storage

import (
	"errors"

	"github.com/Sadere/ya-metrics/internal/common"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

// Интерфейс для хранения данных о метриках
type MetricRepository interface {
	Get(common.MetricType, string) (common.Metrics, error)
	Set(common.Metrics) error

	GetData() (map[string]common.Metrics, error)
	SetData(map[string]common.Metrics) error
}
