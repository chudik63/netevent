package main

import (
	"context"
	"os"

	_ "github.com/lib/pq"

	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/config"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/internal/application/notification"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/logger"
)

func main() {
	ctx := context.Background()

	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) (exitCode int) {
	app := notification.New()

	cfg, err := config.New()
	if err != nil {
		logger.Default().Errorf(ctx, "failed to get config: %s", err)
		return 1
	}

	if err := app.Initialize(ctx, cfg); err != nil {
		logger.Default().Errorf(ctx, "failed to initialize app: %s", err)
		return 1
	}

	if err := app.Run(ctx); err != nil {
		logger.Default().Errorf(ctx, "failed to run app: %s", err)
		return 1
	}

	return 0
}
