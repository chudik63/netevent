package grpc

import (
	"context"
	"net"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/config"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/api/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(ctx context.Context, cfg *config.Config, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCServerPort)

	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	authClient := NewClient(cfg)

	event.RegisterEventServiceServer(grpcServer, NewEventService(ctx, authClient, service))

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
