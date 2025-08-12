package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sayyidinside/gofiber-clean-fresh/shared/config"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func RateLimiter() func(*fiber.Ctx) error {
	cfg := config.AppConfig

	return limiter.New(limiter.Config{
		Max:               cfg.RateLimitMax,
		Expiration:        time.Duration(cfg.RateLimitExp) * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusTooManyRequests,
				Success: false,
				Message: "Limit reached, too many requests detected in short time",
			})
		},
	})
}
