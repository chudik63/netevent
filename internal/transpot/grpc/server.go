package grpc

import (
	"context"
	"event_service/internal/logger"
	"net"

	"event_service/pkg/api/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context, port string, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
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
