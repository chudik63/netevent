package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/config"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/repository"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/service"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/transpot/grpc"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/migrator"

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

	migrator.Start(ctx, cfg)

	eventRepository := repository.New(db)
	eventService := service.New(eventRepository)

	grpcServer, err := grpc.NewServer(ctx, cfg.GRPCServerPort, eventService)
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
