package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	client "gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/client"
	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/config"
	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/transport/grpc"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/logger"
	"go.uber.org/zap"
)

const (
	serviceName = "api-gateway"
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

	authClient := client.NewAuthClient(ctx, cfg)
	eventClient := client.NewEventClient(ctx, cfg)

	grpcserver, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort)
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
	mainLogger.Info(ctx, "server stopped")
}
