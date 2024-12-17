package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/server"
	logger "gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/loger"
	"go.uber.org/zap"
)

var (
	srvGrpcPort = "5100"
)

func main() {
	mainLog := logger.New(logger.ServiceName)
	mainLog.Info("start auth service")

	db, err := postgres.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.ServiceName, mainLog)
	srv := server.New(ctx, srvGrpcPort, db)
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
