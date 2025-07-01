package grpc

import (
	"net"

	"github.com/Sadere/ya-metrics/internal/server/config"
	"github.com/Sadere/ya-metrics/internal/server/grpc/interceptors"
	"github.com/Sadere/ya-metrics/internal/server/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/Sadere/ya-metrics/proto/metrics/v1"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceV1Server

	config        config.Config
	metricService *service.MetricService
	log           *zap.Logger
}

func NewServer(cfg config.Config, mServ *service.MetricService, log *zap.Logger) *MetricsServer {
	return &MetricsServer{
		config:        cfg,
		metricService: mServ,
		log:           log,
	}
}

func (s *MetricsServer) Register() (*grpc.Server, error) {
	srvInterceptors := make([]grpc.UnaryServerInterceptor, 0)

	// Логи
	srvInterceptors = append(srvInterceptors, interceptors.Logger(s.log))

	// Доверенная подсеть
	if len(s.config.TrustedSubnet) > 0 {
		_, trustedSubnet, err := net.ParseCIDR(s.config.TrustedSubnet)

		if err != nil {
			return nil, err
		}

		srvInterceptors = append(srvInterceptors, interceptors.ValidateIP(trustedSubnet))
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		srvInterceptors...,
	))

	// регистрируем сервис
	pb.RegisterMetricsServiceV1Server(srv, s)

	return srv, nil
}
