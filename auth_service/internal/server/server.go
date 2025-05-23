package server

import (
	"context"
	"fmt"
	"net"

	"github.com/chudik63/netevent/auth_service/internal/config"
	"github.com/chudik63/netevent/auth_service/internal/db/postgres"
	"github.com/chudik63/netevent/auth_service/internal/db/postgres/repository"
	pb "github.com/chudik63/netevent/auth_service/pkg/proto"
	"github.com/chudik63/netevent/events_service/pkg/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listen     net.Listener
	db         *postgres.DB
}

func New(ctx context.Context, cfg *config.Config, db *postgres.DB, repo *repository.UserRepository) *Server {
	srvLogger := logger.GetLoggerFromCtx(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		srvLogger.Fatal(ctx, "failed to listen", zap.String("err: ", err.Error()))
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptorLogger(srvLogger)),
	}
	s := grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(s, &Auth{
		repo:        repo,
		eventAdress: cfg.EventsServiceHost + ":" + cfg.EventsServicePort,
	})
	srvLogger.Info(ctx, "server listening at", zap.Int("port", lis.Addr().(*net.TCPAddr).Port))

	return &Server{s, lis, db}
}

func (s *Server) Start(ctx context.Context) error {
	return s.grpcServer.Serve(s.listen)
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	l.Info(ctx, "grpc server Stoped!")
}
