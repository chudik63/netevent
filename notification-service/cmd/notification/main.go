package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/notification"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

const stopTimeout = 3 * time.Second

func main() {
	lg := logger.New(os.Stdout, slog.LevelInfo, "notification-service")
	ctx := logger.CtxWithLogger(context.Background(), lg)

	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) (exitCode int) {
	app := notification.New()

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to get config: %s", err)
		return 1
	}

	if err := app.Initialize(ctx, cfg); err != nil {
		logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to initialize app: %s", err)
		return 1
	}

	go func() {
		if err := app.Run(ctx); err != nil {
			logger.GetLoggerFromCtx(ctx).Errorf(ctx, "failed to run app: %s", err)
			os.Exit(2)
		}
	}()

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)
	<-graceCh

	ctx, cancel := context.WithTimeout(ctx, stopTimeout)
	defer cancel()

	logger.GetLoggerFromCtx(ctx).Infof(ctx, "Graceful shutdown...")

	if err := app.Stop(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Infof(ctx, "Telegram application shutdown error: %s", err)
		return 3
	}

	return 0
}
