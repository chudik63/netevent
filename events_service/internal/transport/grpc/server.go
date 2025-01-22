package grpc

import (
	"context"
	"net"

	"github.com/chudik63/netevent/events_service/internal/config"
	"github.com/chudik63/netevent/events_service/pkg/api/proto/event"
	"github.com/chudik63/netevent/events_service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(ctx context.Context, cfg *config.Config, service Service) (*Server, error) {
	l := logger.GetLoggerFromCtx(ctx)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCServerPort)

	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptor()),
	}
	s := grpc.NewServer(opts...)
	event.RegisterEventServiceServer(grpcServer, NewEventService(ctx, service))

	reflection.Register(grpcServer)

	return &Server{grpcServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	l := logger.GetLoggerFromCtx(ctx)
	l.Info(ctx, "starting grpc server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
