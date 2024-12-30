package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/server"
	loger "github.com/chudik63/netevent/auth/pkg/logger"

	"go.uber.org/zap"
)

var (
	srvGrpcPort = "5100"
)

func main() {
	mainLog := loger.New(loger.ServiceName)
	mainLog.Info("start auth service")

	db, err := postgres.New()
	if err != nil {
		panic(err)
	}
	err = postgres.StartMigration(db.Db.DB)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, loger.ServiceName, mainLog)
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
