package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServerPort string `env:"GRPC_SERVER_PORT" env-default:"9090"`
}

func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("./configs/config.env", &cfg)

	if err != nil {
		return nil
	}

	return &cfg
}
