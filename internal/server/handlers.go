package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) updateHandle(c *gin.Context) {
	metricType := c.Param("type")

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

	err = s.repository.Set(common.Metrics{
		ID:    name,
		MType: string(common.GaugeMetric),
		Value: &valueFloat,
	})
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
}

func (s *Server) updateCounterHandle(c *gin.Context) {
	name := c.Param("metric")
	value := c.Param("value")

	addValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = s.addOrSetCounter(name, addValue)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
}

func (s *Server) addOrSetCounter(name string, addValue int64) (common.Metrics, error) {
	metric, err := s.repository.Get(common.CounterMetric, name)
	if err != nil {
		// Создаем новую метрику если нет такой
		deltaVar := int64(0)
		metric = common.Metrics{
			ID:    name,
			MType: string(common.CounterMetric),
			Delta: &deltaVar,
		}
	}

	*metric.Delta += addValue

	err = s.repository.Set(metric)
	if err != nil {
		return metric, err
	}

	return metric, nil
}

func (s *Server) getMetricHandle(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("metric")

	switch metricType {
	case string(common.CounterMetric):
		metric, err := s.repository.Get(common.CounterMetric, metricName)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		resultDelta := strconv.FormatInt(*metric.Delta, 10)

		c.String(http.StatusOK, resultDelta)
	case string(common.GaugeMetric):
		metric, err := s.repository.Get(common.GaugeMetric, metricName)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		resultValue := strconv.FormatFloat(*metric.Value, 'f', 2, 64)

		c.String(http.StatusOK, resultValue)
	default:
		c.String(http.StatusNotFound, "unknown metric type")
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
}

func (s *Server) getAllMetricsHandle(c *gin.Context) {
	type metric struct {
		Name  string
		Value string
	}

	data, err := s.repository.GetData()
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
