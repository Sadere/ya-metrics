package rest

import (
	"errors"
	"net/http"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/database"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/gin-gonic/gin"
)

func (s *Server) updateHandleJSON(c *gin.Context) {
	var metric common.Metrics
	var err error

	if err = c.BindJSON(&metric); err != nil {
		c.String(http.StatusBadRequest, "failed to parse provided metric")
		return
	}

	metric, err = s.metricService.UpdateMetric(metric)

	if errors.Is(err, database.ErrDBConnection) {
		c.String(http.StatusInternalServerError, "server's storage is down")
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

	if err = c.BindJSON(&metric); err != nil {
		c.String(http.StatusBadRequest, "failed to parse provided metric")
		return
	}

	metric, err = s.metricService.GetMetric(common.MetricType(metric.MType), metric.ID)

	if errors.Is(err, service.ErrWrongMetricType) {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, database.ErrDBConnection) {
		c.String(http.StatusInternalServerError, "server's storage is down")
		return
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

	if err = c.BindJSON(&metrics); err != nil {
		msg := "failed to parse metric list"

		s.log.Sugar().Error(msg)

		c.String(http.StatusBadRequest, msg)
		return
	}

	for _, metric := range metrics {
		_, err = s.metricService.UpdateMetric(metric)

		if errors.Is(err, service.ErrWrongMetricType) {
			s.log.Sugar().Error(err.Error())

			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, database.ErrDBConnection) {
			s.log.Sugar().Error(err.Error())

			c.String(http.StatusInternalServerError, "server's storage is down")
			return
		}

		if err != nil {
			s.log.Sugar().Error(err.Error())

			c.String(http.StatusInternalServerError, "unknown error")
			return
		}
	}

	c.Status(http.StatusOK)
}
