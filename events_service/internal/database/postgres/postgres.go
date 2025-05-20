package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chudik63/netevent/events_service/internal/config"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(ctx context.Context, config config.PostgresConfig) (DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DBName, config.Host, config.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return DB{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return DB{db}, nil
}
