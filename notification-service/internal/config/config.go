package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"gitlab.crja72.ru/gospec/go9/netevent/notification-service/pkg/mail/gmail"
)

type Config struct {
	gmail.GmailConfig
}

func New() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadConfig("./.env", &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &cfg, nil
}
