package main

import (
	"time"

	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/cmd/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/cmd/worker"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/shutdown"
)

func main() {
	depedency := bootstrap.InitApp()

	handler := bootstrap.Initialize(depedency.DB, depedency.Redis.CacheClient, depedency.Redis.LockClient, depedency.RabbitMQ)

	shutdownHandler := shutdown.NewHandler(nil, depedency.DB, depedency.Redis, depedency.RabbitMQ).WithTimeout(30 * time.Second)

	worker.Consume(depedency, shutdownHandler, handler)
}
