package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServerPort string `env:"GRPC_SERVER_PORT"`
	RestServerPort string `env:"REST_SERVER_PORT"`

	AuthServiceHost string `env:"AUTH_SERVICE_HOST"`
	AuthServicePort string `env:"AUTH_SERVICE_PORT"`

	EventServiceHost string `env:"EVENTS_SERVICE_HOST"`
	EventServicePort string `env:"EVENTS_SERVICE_PORT"`
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
