package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string, ttl time.Duration)
	Delete(ctx context.Context, key string)
	Close(ctx context.Context) error
}

func WithCache[T any](ctx context.Context, cache Cache, key string, ttl time.Duration, exec func() (*T, error)) (*T, error) {
	slog.InfoContext(ctx, "getting cache", "key", key)
	value, found := cache.Get(ctx, key)
	var result *T
	if !found {
		slog.InfoContext(ctx, "cache miss", "key", key)
		value, err := exec()
		if err != nil {
			return result, err
		}
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return result, err
		}
		cache.Set(ctx, key, string(jsonValue), ttl)
		return value, nil
	}
	slog.InfoContext(ctx, "cache hit", "key", key)
	var val T
	if err := json.Unmarshal([]byte(value), &val); err != nil {
		return result, err
	}
	return &val, nil
}

func WithRefreshCache[T any](ctx context.Context, cache Cache, key string, ttl time.Duration, value *T) (*T, error) {
	slog.InfoContext(ctx, "refreshing cache", "key", key)
	_, found := cache.Get(ctx, key)
	if found {
		slog.InfoContext(ctx, "deleting cache", "key", key)
		cache.Delete(ctx, key)
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "setting cache", "key", key)
	cache.Set(ctx, key, string(jsonValue), ttl)
	return value, nil
}

func WithDeleteCache(ctx context.Context, cache Cache, key string) error {
	_, found := cache.Get(ctx, key)
	if found {
		slog.InfoContext(ctx, "deleting cache", "key", key)
		cache.Delete(ctx, key)
	}
	return nil
}
