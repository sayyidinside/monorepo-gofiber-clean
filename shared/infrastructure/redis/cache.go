package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
)

type CacheClient struct {
	client *redis.Client
}

func NewCacheClient(cfg *config.Config) *CacheClient {
	return &CacheClient{
		client: redis.NewClient(
			&redis.Options{
				Addr:     cfg.RedisAddress,
				Password: cfg.RedisPassword,
				DB:       0,
			},
		),
	}
}

func (c *CacheClient) Get(ctx context.Context, key string, dest any) (data string, err error) {
	data, err = c.client.Get(ctx, key).Result()

	return
}

func (c *CacheClient) GetObject(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (c *CacheClient) Exist(ctx context.Context, key string) (exist bool, err error) {
	data, err := c.client.Exists(ctx, key).Result()

	if data == 1 {
		exist = true
	}

	return exist, err
}

func (c *CacheClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	var data interface{}

	// Only serialize non-primitive types
	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		data = value
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(jsonData)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *CacheClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *CacheClient) Shutdown() error {
	return c.client.Close()
}

func (c *CacheClient) HealthCheck(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return err
	}
	return c.client.Ping(ctx).Err()
}
