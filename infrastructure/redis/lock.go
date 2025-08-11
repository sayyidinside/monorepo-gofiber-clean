package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

type LockClient struct {
	client *redis.Client
}

func NewLockClient(cfg *config.Config) *LockClient {
	return &LockClient{
		client: redis.NewClient(
			&redis.Options{
				Addr:     cfg.RedisAddress,
				Password: cfg.RedisPassword,
				DB:       1,
			},
		),
	}
}

func (l *LockClient) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return l.client.SetNX(ctx, key, "locked", ttl).Result()
}

func (l *LockClient) ReleaseLock(ctx context.Context, key string) error {
	return l.client.Del(ctx, key).Err()
}

func (l *LockClient) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	acquired, err := l.AcquireLock(ctx, key, ttl)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("failed to acquire lock")
	}

	defer l.ReleaseLock(ctx, key)
	return fn()
}

func (l *LockClient) Shutdown() error {
	return l.client.Close()
}
