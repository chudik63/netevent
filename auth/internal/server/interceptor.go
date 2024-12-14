package server

import (
	"context"

	logger "gitlab.crja72.ru/gospec/go9/netevent/auth_service/pkg/loger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func interceptorLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info("request started", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}
