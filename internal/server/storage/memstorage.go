package storage

import (
	"errors"
	"fmt"

	"github.com/Sadere/ya-metrics/internal/common"
)

// Хранение данных метрик в памяти
type MemMetricRepository struct {
	MetricCounters map[string]int64
	MetricGauges   map[string]float64
}

func NewMemRepository() *MemMetricRepository {
	return &MemMetricRepository{
		MetricCounters: make(map[string]int64),
		MetricGauges:   make(map[string]float64),
	}
}

func (m MemMetricRepository) Get(metricType common.MetricType, key string) (common.Metric, error) {
	metric := common.Metric{}

	if len(key) == 0 {
		return metric, errors.New("key shouldn't be empty")
	}

	switch metricType {
	case common.CounterMetric:
		value, ok := m.MetricCounters[key]
		if !ok {
			return metric, fmt.Errorf("no data with %s key", key)
		}
		metric.Type = common.CounterMetric
		metric.ValueCounter = value
	case common.GaugeMetric:
		value, ok := m.MetricGauges[key]
		if !ok {
			return metric, fmt.Errorf("no data with %s key", key)
		}
		metric.Type = common.GaugeMetric
		metric.ValueGauge = value
	default:
		return metric, fmt.Errorf("invalid metric type %s ", metricType)
	}

	return metric, nil
}

func (m MemMetricRepository) Set(key string, metric common.Metric) error {
	if len(key) == 0 {
		return errors.New("key shouldn't be empty")
	}

	switch metric.Type {
	case common.CounterMetric:
		m.MetricCounters[key] = metric.ValueCounter
	case common.GaugeMetric:
		m.MetricGauges[key] = metric.ValueGauge
	default:
		return fmt.Errorf("invalid metric type %s ", metric.Type)
	}

	return nil
}

func (m MemMetricRepository) GetData() map[string]common.Metric {
	result := make(map[string]common.Metric)

	for k, v := range m.MetricCounters {
		result[k] = common.Metric{
			Type:         common.CounterMetric,
			ValueCounter: v,
		}
	}

	for k, v := range m.MetricGauges {
		result[k] = common.Metric{
			Type:       common.CounterMetric,
			ValueGauge: v,
		}
	}

	return result
}
