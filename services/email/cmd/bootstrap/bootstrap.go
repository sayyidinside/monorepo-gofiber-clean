package bootstrap

import (
	"log"

	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/interfaces/broker/handler"
	sharedRepository "github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	sharedBootstrap "github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/rabbitmq"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB, cacheRedis *redis.CacheClient, lockRedis *redis.LockClient, rabbitMQ *rabbitmq.RabbitMQClient) *handler.Handlers {
	// Repositories
	userRepo := sharedRepository.NewUserRepository(db)

	// Service
	emailService := service.NewEmailService(cacheRedis, userRepo)

	// Handler
	emailHandler := handler.NewEmailHandler(emailService)

	// Setup handler to send to routes setup
	handler := &handler.Handlers{
		EmailHandler: emailHandler,
	}

	return handler
}

func InitApp() *sharedBootstrap.Deps {
	deps, err := sharedBootstrap.NewDeps("../../.env", true)
	if err != nil {
		log.Fatalf("Failed to connect to depedency: %v", err)
	}

	return deps
}
