package config

import (
	"errors"

	"github.com/chudik63/netevent/auth_service/internal/db/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	GRPCServerPort    string `env:"AUTH_SERVICE_PORT"`
	EventsServicePort string `env:"EVENTS_SERVICE_PORT"`
	EventsServiceHost string `env:"EVENTS_SERVICE_HOST"`
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
