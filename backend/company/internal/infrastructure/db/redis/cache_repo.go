package infra_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo[KeyT, ValT any] struct {
	rdb    *redis.Client
	prefix string
}

func NewRepo[KeyT, ValT any](client *redis.Client, prefix string) *CacheRepo[KeyT, ValT] {
	return &CacheRepo[KeyT, ValT]{
		rdb:    client,
		prefix: prefix,
	}
}

func (repo *CacheRepo[KeyT, ValT]) Get(ctx context.Context, key KeyT) *ValT {
	ctx, cancel := context.WithTimeout(ctx, 80*time.Millisecond)
	defer cancel()

	rKey := repo.key(key)
	data, err := repo.rdb.Get(ctx, rKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		// if we have redis issue, we treat it as cache miss - instead of finishing api req with 500
		slog.Warn("redis get error", "err", err)
		return nil
	}

	var val ValT
	err = json.Unmarshal(data, &val)
	if err != nil {
		slog.Warn("couldn't unmarshal: ", repo.prefix, "redis repo:", err)
		return nil
	}

	return &val
}

func (repo *CacheRepo[KeyT, ValT]) Put(ctx context.Context, key KeyT, val *ValT, exp time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	data, mErr := json.Marshal(val)
	if mErr != nil {
		slog.Warn("couldn't marshal in", repo.prefix, "redis repo:", mErr)
		return
	}

	rKey := repo.key(key)
	err := repo.rdb.Set(ctx, rKey, data, exp).Err()

	if err != nil {
		slog.Warn("redis set error", "err", err)
	}
}

func (repo *CacheRepo[KeyT, ValT]) Del(ctx context.Context, key KeyT) {
	ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	rKey := repo.key(key)
	err := repo.rdb.Del(ctx, rKey).Err()

	if err != nil {
		slog.Warn("redis del error", "err", err)
	}
}

func (repo *CacheRepo[KeyT, ValT]) key(key KeyT) string {
	return repo.prefix + fmt.Sprint(key)
}
