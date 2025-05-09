package migrator

import (
	"fmt"

	"github.com/chudik63/netevent/events_service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Start(cfg *config.Config) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Config.UserName, cfg.Config.Password, cfg.Config.Host, cfg.Config.Port, cfg.Config.DbName)

	m, err := migrate.New("file://"+cfg.MigrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to make migration up: %w", err)
	}

	return nil
}
