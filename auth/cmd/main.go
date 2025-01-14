package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chudik63/netevent/auth/internal/config"
	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/server"
	"github.com/chudik63/netevent/event_service/pkg/logger"
	"go.uber.org/zap"
)

const (
	serviceName = "auth"
)

func main() {
	mainLog, err := logger.New(serviceName)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %v", err))
	}

	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLog)

	cfg, err := config.New()
	if err != nil {
		mainLog.Fatal(ctx, "failed to read config", zap.String("err", err.Error()))
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		mainLog.Fatal(ctx, "failed to create database", zap.String("err", err.Error()))
	}
	err = postgres.StartMigration(db.Db.DB)
	if err != nil {
		mainLog.Fatal(ctx, "failed to start migration", zap.String("err", err.Error()))
	}

	srv := server.New(ctx, cfg.GRPCServerPort, db)
	go func() {
		if err := srv.Start(ctx); err != nil {
			mainLog.Error(ctx, "err server", zap.Error(err))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	go func() {
		srv.Stop(ctx)
	}()

}
