package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/db/postgres/repository"
	"github.com/chudik63/netevent/auth/pkg/logger"
	pb "github.com/chudik63/netevent/auth/pkg/proto"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listen     net.Listener
	db         *postgres.DB
}

func New(ctx context.Context, port string, db *postgres.DB) *Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srvLogger := logger.CtxGetLogger(ctx)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptorLogger(srvLogger)),
	}
	s := grpc.NewServer(opts...)

	repo := repository.NewUserRepository(db)
	pb.RegisterAuthServiceServer(s, &Auth{repo: repo})
	log.Printf("server listening at %v", lis.Addr())

	return &Server{s, lis, db}
}

func (s *Server) Start(ctx context.Context) error {
	return s.grpcServer.Serve(s.listen)
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.GracefulStop()
	l := logger.CtxGetLogger(ctx)
	l.Info("grpc server Stoped!")
}
