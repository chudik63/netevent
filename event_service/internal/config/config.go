package config

import (
	"errors"

	"github.com/chudik63/netevent/event_service/internal/database/cache"
	"github.com/chudik63/netevent/event_service/internal/database/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	cache.RedisConfig
	MigrationsPath string `env:"EVENTS_MIGRATIONS_PATH"`
	GRPCServerPort string `env:"EVENTS_SERVICE_PORT"`

	KafkaHost string `env:"KAFKA_HOST"`
	KafkaPort string `env:"KAFKA_PORT"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)

	if cfg == (Config{}) {
		return nil, errors.New("config is empty")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
