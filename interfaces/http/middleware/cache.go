package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

func Cache() func(*fiber.Ctx) error {
	cfg := config.AppConfig

	return cache.New(cache.Config{
		CacheControl: true,
		Expiration:   time.Duration(cfg.CacheExp) * time.Second,
	})
}
