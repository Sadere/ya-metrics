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

	if name == "" || value == "" {
		c.String(http.StatusNotFound, "Insufficient parameters")
		return
	}

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = s.repository.Set(name, common.Metric{
		Type:       common.GaugeMetric,
		ValueGauge: valueFloat,
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

	if name == "" || value == "" {
		c.String(http.StatusNotFound, "Insufficient parameters")
		return
	}

	metric, err := s.repository.Get(common.CounterMetric, name)
	if err != nil {
		// Создаем новую метрику если нет такой
		metric = common.Metric{
			Type: common.CounterMetric,
		}
	}

	addValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	metric.ValueCounter += addValue

	err = s.repository.Set(name, metric)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")

	fmt.Println(s.repository)
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
		c.String(http.StatusOK, fmt.Sprintf("%v", metric.ValueCounter))
	case string(common.GaugeMetric):
		metric, err := s.repository.Get(common.GaugeMetric, metricName)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%v", metric.ValueGauge))
	default:
		c.String(http.StatusNotFound, "unknown metric type")
	}
}

func (s *Server) getAllMetricsHandle(c *gin.Context) {
	type metric struct {
		Name  string
		Value string
	}

	data := s.repository.GetData()
	metrics := make([]metric, len(data))

	for k, v := range data {
		metrics = append(metrics, metric{
			Name:  k,
			Value: fmt.Sprintf("%v", v),
		})
	}

	fmt.Println(metrics)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Metrics": metrics,
	})
}
