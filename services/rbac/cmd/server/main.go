package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/rbac/cmd/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/shutdown"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

func main() {
	depedency := bootstrap.InitApp()

	app := fiber.New(fiber.Config{
		AppName:                 depedency.Config.AppName,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
	})

	// Initialize default config
	app.Use(logger.New())

	// Add Request ID middleware
	app.Use(requestid.New())

	app.Use(helpers.APILogger(helpers.GetAPILogger()))

	// Recover panic
	app.Use(helpers.RecoverWithLog())

	app.Use(helpers.ErrorHelper)

	bootstrap.Initialize(app, depedency.DB, depedency.Redis.CacheClient, depedency.Redis.LockClient)

	app.Use(helpers.NotFoundHelper)

	shutdownHandler := shutdown.NewHandler(app, depedency.DB, depedency.Redis).WithTimeout(30 * time.Second)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", depedency.Config.APort)); err != nil {
			log.Panic(err)
		}
	}()

	shutdownHandler.Listen()
}
