package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Логирует все запросы к серверу
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Запрос
		c.Next()

		// Считаем продолжительность запроса в мили секундах
		latency := time.Since(t)
		miliSeconds := fmt.Sprintf("%d ms", latency.Milliseconds())

		// Код статуса и размер тела запроса
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Тело лога
		logParams := []interface{}{
			"method", c.Request.Method,
			"uri", c.Request.URL,
			"status", status,
			"duration", miliSeconds,
			"size", size,
		}

		// Пишем в лог
		if status >= 500 {
			logger.Sugar().Errorln(logParams...)
		} else {
			logger.Sugar().Infoln(logParams...)
		}
	}
}
