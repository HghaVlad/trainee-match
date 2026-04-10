package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *Config) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.Addr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		PoolTimeout:  500 * time.Millisecond,
	}

	rdb := redis.NewClient(opts)
	err := ping(rdb, time.Second*5)
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func ping(rd *redis.Client, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return rd.Ping(ctx).Err()
}
