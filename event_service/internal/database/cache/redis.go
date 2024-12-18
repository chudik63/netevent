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
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port string `env:"REDIS_PORT" env-default:"6379"`
}

type Cache struct {
	*redis.Client
}

func New(cfg RedisConfig) Cache {
	return Cache{redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})}
}
