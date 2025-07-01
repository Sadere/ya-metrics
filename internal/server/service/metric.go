package service

import (
	"errors"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/storage"
)

var (
	ErrWrongMetricType = errors.New("unknown metric type")
)

type MetricService struct {
	repository storage.MetricRepository
}

func NewMetricService(rep storage.MetricRepository) *MetricService {
	return &MetricService{
		repository: rep,
	}
}

func (s *MetricService) UpdateMetric(metric common.Metrics) (common.Metrics, error) {
	var err error

	switch metric.MType {
	case string(common.GaugeMetric):
		metric, err = s.updateGauge(metric)
	case string(common.CounterMetric):
		metric, err = s.updateCounter(metric)
	default:
		return metric, ErrWrongMetricType
	}

	return metric, err
}

func (s *MetricService) updateGauge(metric common.Metrics) (common.Metrics, error) {
	return metric, s.repository.Set(metric)
}

func (s *MetricService) updateCounter(metric common.Metrics) (common.Metrics, error) {
	metricOld, err := s.repository.Get(common.CounterMetric, metric.ID)
	if err == nil {
		*metric.Delta += *metricOld.Delta
	}

	return metric, s.repository.Set(metric)
}

func (s *MetricService) GetMetric(metricType common.MetricType, ID string) (common.Metrics, error) {
	if metricType != common.CounterMetric && metricType != common.GaugeMetric {
		return common.Metrics{}, ErrWrongMetricType
	}

	return s.repository.Get(metricType, ID)
}


func (s *MetricService) GetAllMetrics() (map[string]common.Metrics, error) {
	return s.repository.GetData()
}

func (s *MetricService) SaveMetrics(metrics map[string]common.Metrics) error {
	return s.repository.SetData(metrics)
}