package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database DB
	Mail     Mail
	Kafka    Kafka
	Sender   Sender
}

type DB struct {
	SQL SQL
}

type SQL struct {
	ConnectionString string `env:"NOTIFICATION_DB_URL"    env-default:"postgres://postgres:123@localhost:5432/netevent?sslmode=disable"`
	Driver           string `env:"NOTIFICATION_DB_DRIVER" env-default:"postgres"`
}

type Mail struct {
	Host     string `env:"MAIL_HOST"`
	Port     int    `env:"MAIL_PORT" env-default:"587"`
	Username string `env:"MAIL_USERNAME"`
	Password string `env:"MAIL_PASSWORD"`
}

type Kafka struct {
	Host      string `env:"KAFKA_HOST"  env-default:"localhost"`
	Port      int    `env:"KAFKA_PORT"  env-default:"9092"`
	Group     string `env:"KAFKA_GROUP" env-default:"registration-group"`
	Topic     string `env:"KAFKA_TOPIC" env-default:"registration"`
	Partition int32  `env:"KAFKA_PARTITION" env-default:"0"`
}

type Sender struct {
	SecondInterval int `env:"SENDER_SECOND_INTERVAL"  env-default:"3600"`
}

func New() (*Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cleanenv.ReadEnv(): %w", err)
	}

	return &cfg, nil
}
