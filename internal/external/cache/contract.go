package cache

import (
	"context"
	"encoding/json"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string, ttl time.Duration)
	Delete(ctx context.Context, key string)
	Close(ctx context.Context) error
}

func WithCache[T any](ctx context.Context, cache Cache, key string, ttl time.Duration, exec func() (*T, error)) (*T, error) {
	value, found := cache.Get(ctx, key)
	var result *T
	if !found {
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
	var val T
	if err := json.Unmarshal([]byte(value), &val); err != nil {
		return result, err
	}
	return &val, nil
}

func WithRefreshCache[T any](ctx context.Context, cache Cache, key string, ttl time.Duration, value *T) (*T, error) {
	_, found := cache.Get(ctx, key)
	if found {
		cache.Delete(ctx, key)
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	cache.Set(ctx, key, string(jsonValue), ttl)
	return value, nil
}

func WithDeleteCache(ctx context.Context, cache Cache, key string) error {
	_, found := cache.Get(ctx, key)
	if found {
		cache.Delete(ctx, key)
	}
	return nil
}
