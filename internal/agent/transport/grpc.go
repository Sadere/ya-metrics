package transport

import (
	"context"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/Sadere/ya-metrics/proto/metrics/v1"
)

// Транспорт отправки метрик используя gRPC
type GRPCMetricTransport struct {
	config config.Config
	client pb.MetricsServiceV1Client
}

func NewGRPCMetricTransport(cfg config.Config) (*GRPCMetricTransport, error) {
	c, err := grpc.NewClient(cfg.ServerAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewMetricsServiceV1Client(c)

	return &GRPCMetricTransport{
		config: cfg,
		client: client,
	}, nil
}

// Отправка метрик на сервер с помощью gRPC
func (t *GRPCMetricTransport) SendMetrics(metrics []common.Metrics) error {
	var req pb.SaveMetricsBatchRequestV1

	pbMetrics := make([]*pb.Metric, len(metrics))

	for i, m := range metrics {
		pbMetrics[i] = common.ProtoFromMetric(&m)
	}

	req.Metrics = pbMetrics

	md := metadata.New(map[string]string{common.IPHeader: t.config.HostAddress})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := t.client.SaveMetricsBatchV1(ctx, &req)

	if err != nil {
		return err
	}

	return nil
}
