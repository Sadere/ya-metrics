package grpc

import (
	"context"
	"errors"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Sadere/ya-metrics/proto/metrics/v1"
)

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequestV1) (*pb.GetMetricResponseV1, error) {
	var response pb.GetMetricResponseV1

	// Валидация метрики
	if err := validateMetric(in.Metrics); err != nil {
		return nil, err
	}

	m := common.MetricFromProto(in.Metrics)

	// Получаем метрику
	metric, err := s.metricService.GetMetric(common.MetricType(m.MType), m.ID)

	// Обработка ошибок
	if errors.Is(err, service.ErrWrongMetricType) {
		return nil, status.Error(codes.InvalidArgument, "wrong metric type")
	}

	if errors.Is(err, storage.ErrMetricNotFound) {
		return nil, status.Error(codes.NotFound, "metric not found")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Возвращаем метрику в ответе
	response.Metric = common.ProtoFromMetric(&metric)

	return &response, nil
}

func (s *MetricsServer) SaveMetricsBatch(ctx context.Context, in *pb.SaveMetricsBatchRequestV1) (*emptypb.Empty, error) {
	var response emptypb.Empty

	for _, m := range in.Metrics {
		// Валидация метрики
		if err := validateMetric(m); err != nil {
			return nil, err
		}

		metric := common.MetricFromProto(m)

		// Обновляем метрику
		_, err := s.metricService.UpdateMetric(metric)

		if errors.Is(err, service.ErrWrongMetricType) {
			return &response, status.Error(codes.InvalidArgument, "wrong metric type")
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &response, nil
}

// Функция возвращает ошибку если валидация метрики неудачная
func validateMetric(m *pb.Metric) error {
	v, err := protovalidate.New()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to init validator: %s", err)
	}

	if err = v.Validate(m); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to validate metrics: %s", err)
	}

	return nil
}
