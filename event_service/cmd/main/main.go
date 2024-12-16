package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/config"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/cache"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/producer"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/repository"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/service"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/transport/grpc"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/migrator"

	"go.uber.org/zap"
)

const (
	serviceName = "event_service"
)

func main() {
	mainLogger, err := logger.New(serviceName)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %v", err))
	}

	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg, err := config.New()
	if err != nil {
		mainLogger.Fatal(ctx, "failed to read config", zap.String("err", err.Error()))
	}

	db := postgres.New(ctx, cfg.Config)
	redis := cache.New(cfg.RedisConfig)

	producer, err := producer.New(ctx, []string{cfg.Kafka1Host + ":" + cfg.Kafka1Port, cfg.Kafka2Host + ":" + cfg.Kafka2Port, cfg.Kafka3Host + ":" + cfg.Kafka3Port})
	if err != nil {
		mainLogger.Fatal(ctx, "failed to create broker", zap.String("err", err.Error()))
	}

	migrator.Start(ctx, cfg)

	eventRepository := repository.New(db)
	eventService := service.New(eventRepository, redis, producer)

	grpcServer, err := grpc.NewServer(ctx, cfg, eventService)
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
