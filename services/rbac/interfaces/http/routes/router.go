package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/sayyidinside/monorepo-gofiber-clean/infrastructure/config"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/middleware"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/routes/tests"
	v1 "github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/routes/v1"
	"github.com/sayyidinside/monorepo-gofiber-clean/pkg/helpers"
)

func Setup(app *fiber.App, handlers *handler.Handlers) {
	cfg := config.AppConfig

	api := app.Group("/api")
	test := app.Group("/tests")

	// Apply middleware for general API routes
	api.Use(helmet.New())
	api.Use(middleware.CORS())
	api.Use(middleware.WhitelistIP())
	api.Use(middleware.RateLimiter())
	api.Use(middleware.Cache())

	v1.RegisterRoutes(api, handlers)
	tests.SetupApiTestRoutes(test)

	app.Get("/", func(c *fiber.Ctx) error {
		log := helpers.CreateLog(app)
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  200,
			Success: true,
			Message: fmt.Sprintf("App is running, with name %s", cfg.AppName),
			Log:     &log,
		})
	})
}
