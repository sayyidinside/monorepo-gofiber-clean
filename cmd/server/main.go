package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/shutdown"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func main() {
	bootstrap.InitApp()

	app := fiber.New(fiber.Config{
		AppName:                 config.AppConfig.AppName,
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

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := redis.Connect(config.AppConfig)

	bootstrap.Initialize(app, db, redisClient.CacheClient, redisClient.LockClient)

	app.Use(helpers.NotFoundHelper)

	shutdownHandler := shutdown.NewHandler(app, db, redisClient).WithTimeout(30 * time.Second)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port)); err != nil {
			log.Panic(err)
		}
	}()

	shutdownHandler.Listen()
}
