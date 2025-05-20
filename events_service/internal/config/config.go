package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	PostgresConfig struct {
		UserName string `env:"EVENTS_POSTGRES_USER"`
		Password string `env:"EVENTS_POSTGRES_PASSWORD"`
		Host     string `env:"EVENTS_POSTGRES_HOST"`
		Port     string `env:"EVENTS_POSTGRES_PORT"`
		DBName   string `env:"EVENTS_POSTGRES_DB"`
		SSLMode  string `env:"EVENTS_POSTGRES_SSL_MODE"`
	}

	KafkaConfig struct {
		KafkaHost string `env:"KAFKA_HOST"`
		KafkaPort string `env:"KAFKA_PORT"`
	}

	RedisConfig struct {
		Host string `env:"EVENTS_REDIS_HOST"`
		Port string `env:"EVENTS_REDIS_PORT"`
	}

	Config struct {
		Postgres PostgresConfig
		Kafka    KafkaConfig
		Redis    RedisConfig

		MigrationsPath string `env:"EVENTS_MIGRATIONS_PATH"`
		GRPCServerPort string `env:"EVENTS_SERVICE_PORT"`
	}
)

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

func (cfg *Config) GetUserName() string {
	return cfg.Postgres.UserName
}

func (cfg *Config) GetPassword() string {
	return cfg.Postgres.Password
}

func (cfg *Config) GetHost() string {
	return cfg.Postgres.Host
}

func (cfg *Config) GetPort() string {
	return cfg.Postgres.Port
}

func (cfg *Config) GetDBName() string {
	return cfg.Postgres.DBName
}

func (cfg *Config) GetSSLMode() string {
	return cfg.Postgres.SSLMode
}

func (cfg *Config) GetMigrationsPath() string {
	return cfg.MigrationsPath
}
