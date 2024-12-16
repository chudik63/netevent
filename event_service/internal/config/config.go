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

	Kafka1Host string `env:"KAFKA1_HOST"`
	Kafka1Port string `env:"KAFKA1_PORT"`

	Kafka2Host string `env:"KAFKA2_HOST"`
	Kafka2Port string `env:"KAFKA2_PORT"`

	Kafka3Host string `env:"KAFKA3_HOST"`
	Kafka3Port string `env:"KAFKA3_PORT"`
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
