package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config interface {
	GetUserName() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetDBName() string
	GetMigrationsPath() string
	GetSSLMode() string
}

func Start(cfg Config) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.GetUserName(), cfg.GetPassword(), cfg.GetHost(), cfg.GetPort(), cfg.GetDBName(), cfg.GetSSLMode())

	m, err := migrate.New("file://"+cfg.GetMigrationsPath(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to make migration up: %w", err)
	}

	return nil
}
