package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

func NewClient(cfg config.Redis) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		PoolTimeout:  500 * time.Millisecond,
	}

	rdb := redis.NewClient(opts)
	err := ping(rdb, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return rdb, nil
}

func ping(rd *redis.Client, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return rd.Ping(ctx).Err()
}
