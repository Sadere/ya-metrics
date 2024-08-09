package interceptors

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Лог gRPC запросов к серверу
func Logger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		t := time.Now()

		// Обработка запроса
		res, err := handler(ctx, req)

		// Считаем продолжительность запроса в мили секундах
		latency := time.Since(t)
		miliSeconds := fmt.Sprintf("%d ms", latency.Milliseconds())

		// Тело лога
		logParams := []interface{}{
			"request", fmt.Sprintf("%v", req),
			"response", fmt.Sprintf("%v", res),
			"method", info.FullMethod,
			"error", err,
			"duration", miliSeconds,
		}

		// Пишем в лог
		if err != nil {
			logger.Sugar().Errorln(logParams...)
		} else {
			logger.Sugar().Infoln(logParams...)
		}

		return res, err
	}
}
