package grpc

import (
	"context"
	"errors"

	"github.com/Sadere/ya-metrics/internal/common"
	pb "github.com/Sadere/ya-metrics/internal/proto"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"github.com/Sadere/ya-metrics/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse

	m := common.MetricFromProto(in.Metrics)

	metric, err := s.metricService.GetMetric(common.MetricType(m.MType), m.ID)

	if errors.Is(err, service.ErrWrongMetricType) {
		return nil, status.Error(codes.InvalidArgument, "wrong metric type")
	}

	if errors.Is(err, storage.ErrMetricNotFound) {
		return nil, status.Error(codes.NotFound, "metric not found")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Metric = common.ProtoFromMetric(&metric)

	return &response, nil
}

func (s *MetricsServer) SaveMetricsBatch(ctx context.Context, in *pb.SaveMetricsBatchRequest) (*pb.SaveMetricsBatchResponse, error) {
	var response pb.SaveMetricsBatchResponse

	for _, m := range in.Metrics {
		metric := common.MetricFromProto(m)

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