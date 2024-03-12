package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MetricType string

const (
	CounterMetric MetricType = "counter"
	GaugeMetric   MetricType = "gauge"
)

func (s *Server) updateHandle(c *gin.Context) {
	metricType := c.Param("type")

	switch metricType {
	case string(GaugeMetric):
		s.updateGaugeHandle(c)
		return
	case string(CounterMetric):
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

	err = s.storage.SetFloat64(name, valueFloat)
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

	oldValue, err := s.storage.GetInt64(name)
	if err != nil {
		oldValue = 0
	}

	addValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = s.storage.SetInt64(name, addValue+oldValue)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")

	fmt.Println(s.storage)
}

func (s *Server) getMetricHandle(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("metric")

	if metricType != "counter" && metricType != "gauge" {
		c.String(http.StatusNotFound, "unknown metric type")
		return
	}

	metricValue, err := s.storage.Get(metricName)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.String(http.StatusOK, metricValue)
}

func (s *Server) getAllMetricsHandle(c *gin.Context) {
	type metric struct {
		Name  string
		Value string
	}

	data := s.storage.GetData()
	metrics := make([]metric, len(data))

	for k, v := range s.storage.GetData() {
		metrics = append(metrics, metric{
			Name:  k,
			Value: v,
		})
	}

	fmt.Println(metrics)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Metrics": metrics,
	})
}
