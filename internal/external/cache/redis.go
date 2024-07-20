package cache

import (
	"context"
	"crypto/tls"
	"log/slog"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(ctx context.Context, config *config.Config) Cache {
	var tlsConfig *tls.Config

	if !config.ApiConfig.IsDevelopment() {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := &redisCache{
		client: redis.NewClient(&redis.Options{
			Addr:         config.CacheConfig.Host,
			Password:     config.CacheConfig.Password,
			DB:           config.CacheConfig.DB,
			DialTimeout:  100 * time.Millisecond,
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
			TLSConfig:    tlsConfig,
		}),
	}

	if err := client.client.Ping(ctx).Err(); err != nil {
		slog.ErrorContext(ctx, "error connecting to redis", "err", err)
	}

	return client
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		slog.DebugContext(ctx, "cache miss", "key", key)
		return "", false
	}

	slog.DebugContext(ctx, "cache hit", "key", key)
	return val, true
}

func (c *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		slog.ErrorContext(ctx, "error setting cache", "key", key, "err", err)
	}
	slog.DebugContext(ctx, "cache set", "key", key)
}

func (c *redisCache) Delete(ctx context.Context, key string) {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		slog.ErrorContext(ctx, "error deleting cache", "key", key, "err", err)
	}
	slog.DebugContext(ctx, "cache delete", "key", key)
}

func (c *redisCache) Close(ctx context.Context) error {
	slog.InfoContext(ctx, "closing redis client")

	if err := c.client.Close(); err != nil {
		slog.ErrorContext(ctx, "error closing redis client", "err", err)
		return err
	}

	slog.InfoContext(ctx, "redis client closed")

	return nil
}
