package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var path = ".env"

type Config struct {
	UserName string `env:"AUTH_POSTGRES_USER"`
	Password string `env:"AUTH_POSTGRES_PASWD"`
	Host     string `env:"AUTH_POSTGRES_HOST"`
	Port     string `env:"AUTH_POSTGRES_PORT"`
	DBname   string `env:"AUTH_POSTGRES_DATABASE"`
}

type DB struct {
	Db *sqlx.DB
}

func New(cfg Config) (*DB, error) {
	dataSorceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.UserName, cfg.Password, cfg.DBname, cfg.Host, cfg.Port)

	db, err := sqlx.Connect("postgres", dataSorceName)
	if err != nil {
		return nil, err
	}

	if _, err := db.Connx(context.Background()); err != nil {
		return nil, err
	}

	return &DB{Db: db}, nil
}
