package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/pkg/api/gateway"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
}

func New(ctx context.Context, port, restPort int) (*Server, error) {
	logs := logger.GetLoggerFromCtx(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logs.Fatal(ctx, "failed to listen", zap.String("err", err.Error()))
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	gateway.RegisterGatewayServer(grpcServer, NewGatewayServer())

	mux := runtime.NewServeMux()
	gateway.RegisterGatewayHandlerServer(context.Background(), mux, NewGatewayServer())
	restSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: mux,
	}
	return &Server{grpcServer, restSrv, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	logs := logger.GetLoggerFromCtx(ctx)

	eg := errgroup.Group{}

	eg.Go(func() error {
		logs.Info(ctx, fmt.Sprintf("gRPC server listening on %s", s.listener.Addr()))
		return s.grpcServer.Serve(s.listener)
	})

	eg.Go(func() error {
		logs.Info(ctx, fmt.Sprintf("Rest server listening on %s", s.restServer.Addr))
		return s.restServer.ListenAndServe()
	})

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.restServer.Shutdown(ctx)
}
