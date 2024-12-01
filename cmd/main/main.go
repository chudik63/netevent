package main

import (
	"context"
	"event_service/internal/config"
	"event_service/internal/logger"
	"event_service/internal/transpot/grpc"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

const (
	serviceName = "event_service"
)

func main() {
	mainLogger := logger.New(serviceName)
	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg := config.New()

	if cfg == nil {
		mainLogger.Fatal(ctx, "failed to read config")
	}

	grpcServer, err := grpc.New(cfg.GRPCServerPort)
	if err != nil {
		mainLogger.Fatal(ctx, "failed to listen", zap.String("err: ", err.Error()))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			mainLogger.Fatal(ctx, "failed to start grpc server")
		}
	}()

	<-quit

	grpcServer.Stop()
	mainLogger.Info(ctx, "Server stopped")
}
