package cache

import (
	"fmt"
	"time"

	"github.com/chudik63/netevent/events_service/internal/config"
	"github.com/redis/go-redis/v9"
)

const (
	Durability = 2 * time.Minute
)

type Cache struct {
	*redis.Client
}

func New(cfg config.RedisConfig) Cache {
	return Cache{redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})}
}
