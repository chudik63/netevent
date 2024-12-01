package grpc

import (
	"context"
	"event_service/internal/logger"
	"net"

	"event_service/pkg/api/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(port string) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	event.RegisterEventServiceServer(grpcServer, NewEventService())

	return &Server{grpcServer, lis}, nil
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
