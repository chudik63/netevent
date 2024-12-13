package cache

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	Durability = 20 * time.Minute
)

type RedisConfig struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port string `env:"REDIS_PORT" env-default:"6379"`
}

func New(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})
}
