package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/chudik63/netevent/auth/internal/config"
	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/server"
	"github.com/chudik63/netevent/auth/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	mainLog := logger.New(logger.ServiceName)
	mainLog.Info("start auth service")

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		panic(err)
	}
	err = postgres.StartMigration(db.Db.DB)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.ServiceName, mainLog)
	srv := server.New(ctx, cfg.GRPCServerPort, db)
	go func() {
		if err := srv.Start(ctx); err != nil {
			mainLog.Error("err server", zap.Error(err))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	go func() {
		srv.Stop(ctx)
	}()

}
