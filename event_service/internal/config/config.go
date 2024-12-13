package config

import (
	"errors"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/cache"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	cache.RedisConfig
	MigrationsPath string `env:"MIGRATIONS_PATH"`
	GRPCServerPort string `env:"GRPC_SERVER_PORT"`

	AuthServerPort string `env:"AUTH_SERVER_PORT"`
	AuthServerHost string `env:"AUTH_SERVER_HOST"`

	KafkaHost string `env:"KAFKA_HOST"`
	KafkaPort string `env:"KAFKA_HOST"`
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
