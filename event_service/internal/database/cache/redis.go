package cache

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	Durability = 2 * time.Minute
)

type RedisConfig struct {
	Host string `env:"EVENTS_REDIS_HOST"`
	Port string `env:"EVENTS_REDIS_PORT"`
}

type Cache struct {
	*redis.Client
}

func New(cfg RedisConfig) Cache {
	return Cache{redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})}
}
