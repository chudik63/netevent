package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/chudik63/netevent/api-gateway/internal/client"
	"github.com/chudik63/netevent/api-gateway/pkg/api/gateway"
	"github.com/chudik63/netevent/event_service/pkg/logger"

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

func New(ctx context.Context, port, restPort string, authClient *client.AuthClient, eventClient *client.EventClient) (*Server, error) {
	logs := logger.GetLoggerFromCtx(ctx)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logs.Fatal(ctx, "failed to listen", zap.String("err", err.Error()))
	}

	interceptor := NewAuthInterceptor(authClient)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	reflection.Register(grpcServer)
	gateway.RegisterGatewayServer(grpcServer, NewGatewayServer(authClient, eventClient))

	mux := runtime.NewServeMux()
	if err := gateway.RegisterGatewayHandlerServer(context.Background(), mux, NewGatewayServer(authClient, eventClient)); err != nil {
		logs.Fatal(ctx, "failed to create gatewate handler", zap.String("err", err.Error()))
		return nil, err
	}
	restSrv := &http.Server{
		Addr:    ":" + restPort,
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
