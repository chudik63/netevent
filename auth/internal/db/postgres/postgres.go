package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var path = ".env"

type Config struct {
	UserName string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASWD" env-default:"12345678"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DBname   string `env:"POSTGRES_DATABASE" env-default:"Users"`
}

type DB struct {
	Db *sqlx.DB
}

func readConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func New() (*DB, error) {
	cfg, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cfg)

	cfg.Port = "5432"
	dataSorceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.UserName, cfg.Password, cfg.DBname, cfg.Host, cfg.Port)
	fmt.Println(dataSorceName)
	db, err := sqlx.Connect("postgres", dataSorceName)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := db.Connx(context.Background()); err != nil {
		return nil, err
	}
	return &DB{Db: db}, nil
}
