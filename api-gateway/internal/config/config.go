package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-default:"5400"`
	RestServerPort int `env:"REST_SERVER_PORT" env-default:"8000"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./.env", &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
