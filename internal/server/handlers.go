package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) updateHandle(c *gin.Context) {
	metricType := c.Param("type")

	c.Header("Content-Type", "text/plain; charset=utf-8")

	switch metricType {
	case string(common.GaugeMetric):
		s.updateGaugeHandle(c)
		return
	case string(common.CounterMetric):
		s.updateCounterHandle(c)
		return
	default:
		c.String(http.StatusBadRequest, "Unknown metric type")
	}
}

func (s *Server) updateGaugeHandle(c *gin.Context) {
	name := c.Param("metric")
	value := c.Param("value")

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = s.metricService.UpdateMetric(common.Metrics{
		ID:    name,
		MType: string(common.GaugeMetric),
		Value: &valueFloat,
	})
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
}

func (s *Server) updateCounterHandle(c *gin.Context) {
	name := c.Param("metric")
	value := c.Param("value")

	addValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = s.metricService.UpdateMetric(common.Metrics{
		ID:    name,
		MType: string(common.CounterMetric),
		Delta: &addValue,
	})
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}

func (s *Server) getMetricHandle(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("metric")

	metric, err := s.metricService.GetMetric(common.MetricType(metricType), metricName)

	if errors.Is(err, storage.ErrMetricNotFound) || errors.Is(err, service.ErrWrongMetricType) {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	switch metricType {
	case string(common.CounterMetric):
		resultDelta := strconv.FormatInt(*metric.Delta, 10)

		c.String(http.StatusOK, resultDelta)
	case string(common.GaugeMetric):
		resultValue := strconv.FormatFloat(*metric.Value, 'f', 6, 64)
		resultValue = strings.TrimRight(resultValue, "0")

		c.String(http.StatusOK, resultValue)
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
}

func (s *Server) getAllMetricsHandle(c *gin.Context) {
	type metric struct {
		Name  string
		Value string
	}

	data, err := s.metricService.GetAllMetrics()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	metrics := make([]metric, len(data))

	for k, v := range data {
		m := metric{
			Name: k,
		}

		if v.MType == string(common.CounterMetric) {
			m.Value = fmt.Sprintf("%d", *v.Delta)
		}

		if v.MType == string(common.GaugeMetric) {
			m.Value = fmt.Sprintf("%f", *v.Value)
		}

		metrics = append(metrics, m)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Metrics": metrics,
	})
}

func (s *Server) pingHandle(c *gin.Context) {
	err := s.db.Ping()

	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}
