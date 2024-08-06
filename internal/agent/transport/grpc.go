package transport

import (
	"context"
	"fmt"

	"github.com/Sadere/ya-metrics/internal/agent/config"
	"github.com/Sadere/ya-metrics/internal/common"
	pb "github.com/Sadere/ya-metrics/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Транспорт отправки метрик используя gRPC
type GRPCMetricTransport struct {
	config config.Config
	client pb.MetricsClient
}

func NewGRPCMetricTransport(cfg config.Config) (*GRPCMetricTransport, error) {
	c, err := grpc.NewClient(cfg.ServerAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewMetricsClient(c)

	return &GRPCMetricTransport{
		config: cfg,
		client: client,
	}, nil
}

// Отправка метрик на сервер с помощью gRPC
func (t *GRPCMetricTransport) SendMetrics(metrics []common.Metrics) error {
	var req pb.SaveMetricsBatchRequest

	pbMetrics := make([]*pb.Metric, len(metrics))

	for i, m := range metrics {
		pbMetrics[i] = common.ProtoFromMetric(&m)
	}

	req.Metrics = pbMetrics

	md := metadata.New(map[string]string{common.IPHeader: t.config.HostAddress})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := t.client.SaveMetricsBatch(ctx, &req)

	if err != nil {
		return err
	}

	if len(res.Error) > 0 {
		return fmt.Errorf("gRPC error: %s", res.Error)
	}

	return nil
}
