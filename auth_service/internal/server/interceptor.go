package server

import (
	"context"

	"github.com/chudik63/netevent/events_service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func interceptorLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info(ctx, "request started", zap.String("method", info.FullMethod))
		resp, err = handler(ctx, req)
		if err != nil {
			l.Error(ctx, "error in response: ", zap.String("error: ", err.Error()))
		}
		return resp, err
	}
}
