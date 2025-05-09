package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func New(ctx context.Context, config Config) (DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DbName, config.Host, config.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return DB{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return DB{db}, nil
}
