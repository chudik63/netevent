package server

import (
	"context"
	"fmt"
	"net"

	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/db/postgres/repository"
	pb "github.com/chudik63/netevent/auth/pkg/proto"
	"github.com/chudik63/netevent/event_service/pkg/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listen     net.Listener
	db         *postgres.DB
}

func New(ctx context.Context, port string, db *postgres.DB) *Server {
	srvLogger := logger.GetLoggerFromCtx(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		srvLogger.Fatal(ctx, "failed to listen", zap.String("err: ", err.Error()))
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptorLogger(srvLogger)),
	}
	s := grpc.NewServer(opts...)

	repo := repository.NewUserRepository(db)
	pb.RegisterAuthServiceServer(s, &Auth{repo: repo})
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
