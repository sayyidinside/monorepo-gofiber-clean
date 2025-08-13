package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
)

func CORS() func(*fiber.Ctx) error {
	cfg := config.AppConfig

	return cors.New(cors.Config{
		AllowOrigins:  cfg.CorsAllowOrigins,
		AllowMethods:  cfg.CorsAllowMethods,
		MaxAge:        cfg.CorsMaxAge,
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	})
}
