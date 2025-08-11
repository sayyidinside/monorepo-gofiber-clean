package redis

import "github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"

type RedisClient struct {
	CacheClient *CacheClient
	LockClient  *LockClient
}

func Connect(cfg *config.Config) *RedisClient {
	return &RedisClient{
		CacheClient: NewCacheClient(cfg),
		LockClient:  NewLockClient(cfg),
	}
}
