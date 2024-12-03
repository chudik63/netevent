package postgres

import (
	"context"
	"database/sql"
	"event_service/internal/logger"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	UserName string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
}

func New(ctx context.Context, config Config) *sql.DB {
	l := logger.GetLoggerFromCtx(ctx)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DbName, config.Host, config.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		l.Fatal(ctx, "failed to open database", zap.String("err", err.Error()))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		l.Error(ctx, "failed to ping database", zap.String("err", err.Error()))
	}

	return db
}
