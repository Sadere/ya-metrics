package interceptors

import (
	"context"
	"net"

	"github.com/Sadere/ya-metrics/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ValidateIP(trustedSubnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Получаем IP из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "failed to retrieve grpc metadata")
		}

		values := md.Get(common.IPHeader)

		if len(values) == 0 {
			return nil, status.Error(codes.PermissionDenied, "no IP passed")
		}

		IPText := values[0]

		IP := net.ParseIP(IPText)

		// Проверяем входит ли IP в доверенные
		if !trustedSubnet.Contains(IP) {
			return nil, status.Error(codes.PermissionDenied, "IP not allowed")
		}

		// Успешная проверка
		return handler(ctx, req)
	}
}
