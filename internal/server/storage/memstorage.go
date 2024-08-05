package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Sadere/ya-metrics/internal/common"
)

// Хранение данных метрик в памяти
type MemMetricRepository struct {
	MetricCounters map[string]int64
	MetricGauges   map[string]float64
	mu             *sync.RWMutex
}

func NewMemRepository() *MemMetricRepository {
	return &MemMetricRepository{
		MetricCounters: make(map[string]int64),
		MetricGauges:   make(map[string]float64),
		mu:             &sync.RWMutex{},
	}
}

func (m MemMetricRepository) Get(metricType common.MetricType, key string) (common.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metric := common.Metrics{}

	if len(key) == 0 {
		return metric, errors.New("key shouldn't be empty")
	}

	switch metricType {
	case common.CounterMetric:
		value, ok := m.MetricCounters[key]
		if !ok {
			return metric, fmt.Errorf("no data with %s key", key)
		}
		metric.MType = string(common.CounterMetric)
		metric.Delta = &value
	case common.GaugeMetric:
		value, ok := m.MetricGauges[key]
		if !ok {
			return metric, fmt.Errorf("no data with %s key", key)
		}
		metric.MType = string(common.GaugeMetric)
		metric.Value = &value
	default:
		return metric, fmt.Errorf("invalid metric type %s ", metricType)
	}

	metric.ID = key

	return metric, nil
}

func (m MemMetricRepository) Set(metric common.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := metric.ID

	if len(key) == 0 {
		return errors.New("key shouldn't be empty")
	}

	switch metric.MType {
	case string(common.CounterMetric):
		m.MetricCounters[key] = *metric.Delta
	case string(common.GaugeMetric):
		m.MetricGauges[key] = *metric.Value
	default:
		return fmt.Errorf("invalid metric type %s ", metric.MType)
	}

	return nil
}

func (m MemMetricRepository) GetData() (map[string]common.Metrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]common.Metrics)

	for k, v := range m.MetricCounters {
		result[k] = common.Metrics{
			ID:    k,
			MType: string(common.CounterMetric),
			Delta: &v,
		}
	}

	for k, v := range m.MetricGauges {
		result[k] = common.Metrics{
			ID:    k,
			MType: string(common.GaugeMetric),
			Value: &v,
		}
	}

	return result, nil
}

func (m MemMetricRepository) SetData(metrics map[string]common.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range metrics {
		switch v.MType {
		case string(common.CounterMetric):
			m.MetricCounters[k] = *v.Delta
		case string(common.GaugeMetric):
			m.MetricGauges[k] = *v.Value
		}
	}

	return nil
}
