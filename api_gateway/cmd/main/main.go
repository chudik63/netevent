package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	client "github.com/chudik63/netevent/api_gateway/internal/client"
	"github.com/chudik63/netevent/api_gateway/internal/config"
	"github.com/chudik63/netevent/api_gateway/internal/transport/grpc"
	"github.com/chudik63/netevent/events_service/pkg/logger"

	"go.uber.org/zap"
)

const (
	serviceName = "api_gateway"
)

func main() {
	mainLogger, err := logger.New(serviceName)
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg, err := config.New()
	if err != nil {
		mainLogger.Fatal(ctx, "failed to create config", zap.String("err", err.Error()))
	}

	// authClient := client.NewAuthClient(ctx, cfg)
	// eventClient := client.NewEventClient(ctx, cfg)

	authClient := &client.AuthClient{}
	eventClient := &client.EventClient{}

	grpcserver, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort, authClient, eventClient)
	if err != nil {
		log.Fatalf("Error creating grpc server: %v", err)
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcserver.Start(ctx); err != nil {
			mainLogger.Info(ctx, "failed to start server", zap.String("err", err.Error()))
		}
	}()

	<-graceCh
	if err := grpcserver.Stop(ctx); err != nil {
		mainLogger.Info(ctx, "failed to stop server", zap.String("err", err.Error()))
	}

	authClient.Stop()
	eventClient.Stop()

	mainLogger.Info(ctx, "server stopped")
}
