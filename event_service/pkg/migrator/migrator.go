package migrator

import (
	"context"
	"fmt"

	"github.com/chudik63/netevent/event_service/internal/config"
	"github.com/chudik63/netevent/event_service/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func Start(ctx context.Context, cfg *config.Config) {
	l := logger.GetLoggerFromCtx(ctx)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Config.UserName, cfg.Config.Password, cfg.Config.Host, cfg.Config.Port, cfg.Config.DbName)

	m, err := migrate.New("file://"+cfg.MigrationsPath, dbURL)
	if err != nil {
		l.Error(ctx, "failed to create migration", zap.String("err", err.Error()))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		l.Error(ctx, "failed to make migration up", zap.String("err", err.Error()))
	}
}
