package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GRPCServerPort string `env:"GRPC_SERVER_PORT"`
	RestServerPort string `env:"REST_SERVER_PORT"`

	AuthServiceHost string `env:"AUTH_SERVICE_HOST"`
	AuthServicePort string `env:"AUTH_SERVICER_PORT"`

	EventServiceHost string `env:"EVENT_SERVICE_HOST"`
	EventServicePort string `env:"EVENT_SERVICE_PORT"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./.env", &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
