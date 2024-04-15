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
		err = s.repository.Set(metric)
	case string(common.CounterMetric):
		metric, err = s.addOrSetCounter(metric.ID, *metric.Delta)
	default:
		c.String(http.StatusBadRequest, "Unknown metric type")
		return
	}

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
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

func (s *Server) updateBatchHandleJSON(c *gin.Context) {
	var metrics []common.Metrics
	var err error

	if err := c.BindJSON(&metrics); err != nil {
		msg := "failed to parse metric list"

		s.log.Sugar().Error(msg)

		c.String(http.StatusBadRequest, msg)
		return
	}

	for _, metric := range metrics {
		switch metric.MType {
		case string(common.GaugeMetric):
			err = s.repository.Set(metric)
		case string(common.CounterMetric):
			metric, err = s.addOrSetCounter(metric.ID, *metric.Delta)
		default:
			msg := "unknown metric type"

			s.log.Sugar().Error(msg)

			c.String(http.StatusBadRequest, msg)
			return
		}

		if err != nil {
			s.log.Sugar().Error(err.Error())

			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
