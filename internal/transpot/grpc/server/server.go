package server

import (
	"context"
	"event_service/internal/logger"
	"event_service/internal/transpot/grpc/service"
	"net"

	"event_service/pkg/api/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context, port string) *Server {
	l := logger.GetLoggerFromCtx(ctx)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		l.Fatal(ctx, "failed to listen", zap.String("err: ", err.Error()))
	}

	grpcServer := grpc.NewServer()
	event.RegisterEventServiceServer(grpcServer, service.NewEventService())

	return &Server{grpcServer, lis}
}

func (s *Server) Start(ctx context.Context) error {
	l := logger.GetLoggerFromCtx(ctx)
	l.Info(ctx, "starting grpc server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)

	l.Info(ctx, "gRPC server stopped")
}
