package server

import (
	"net/http"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) updateHandleJSON(c *gin.Context) {
	var metric common.Metrics
	var err error

	if err := c.BindJSON(&metric); err != nil {
		c.String(http.StatusBadRequest, "failed to parse provided metric")
		return
	}

	switch metric.MType {
	case string(common.GaugeMetric):
		s.repository.Set(metric.ID, metric)
	case string(common.CounterMetric):
		metric, err = s.addOrSetCounter(metric.ID, *metric.Delta)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
	default:
		c.String(http.StatusBadRequest, "Unknown metric type")
		return
	}

	c.JSON(http.StatusOK, metric)
}

func (s *Server) getMetricHandleJSON(c *gin.Context) {
	var metric common.Metrics
	var err error

	if err := c.BindJSON(&metric); err != nil {
		c.String(http.StatusBadRequest, "failed to parse provided metric")
		return
	}

	switch metric.MType {
	case string(common.GaugeMetric):
		metric, err = s.repository.Get(common.GaugeMetric, metric.ID)
	case string(common.CounterMetric):
		metric, err = s.repository.Get(common.CounterMetric, metric.ID)
	default:
		c.String(http.StatusBadRequest, "Unknown metric type")
	}

	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, metric)
}
