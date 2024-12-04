package main

import (
	"context"
	"event_service/internal/config"
	"event_service/internal/database/postgres"
	"event_service/internal/logger"
	"event_service/internal/repository"
	"event_service/internal/service"
	"event_service/internal/transpot/grpc"
	"event_service/pkg/migrator"
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

	cfg, err := config.New()

	if err != nil {
		mainLogger.Fatal(ctx, "failed to read config", zap.String("err", err.Error()))
	}

	db := postgres.New(ctx, cfg.Config)
	_ = db

	migrator.Start(ctx, cfg)

	eventRepository := repository.New(db)
	eventService := service.New(eventRepository)

	grpcServer, err := grpc.New(cfg.GRPCServerPort, eventService)
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
