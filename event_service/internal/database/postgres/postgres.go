package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	UserName string `env:"EVENTS_POSTGRES_USER"`
	Password string `env:"EVENTS_POSTGRES_PASSWORD"`
	Host     string `env:"EVENTS_POSTGRES_HOST"`
	Port     string `env:"EVENTS_POSTGRES_PORT"`
	DbName   string `env:"EVENTS_POSTGRES_DB"`
}

type DB struct {
	*sql.DB
}

func New(ctx context.Context, config Config) DB {
	l := logger.GetLoggerFromCtx(ctx)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DbName, config.Host, config.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		l.Fatal(ctx, "failed to open database", zap.String("err", err.Error()))
	}

	if err := db.Ping(); err != nil {
		l.Fatal(ctx, "failed to ping database", zap.String("err", err.Error()))
	}

	return DB{db}
}
