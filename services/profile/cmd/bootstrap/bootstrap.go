package bootstrap

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/cmd/worker"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/interfaces/http/routes"
	sharedRepo "github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	sharedBootstrap "github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/http/middleware"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"gorm.io/gorm"
)

func Initialize(app *fiber.App, db *gorm.DB, cacheRedis *redis.CacheClient, lockRedis *redis.LockClient) {
	// Repositories
	profileRepo := repository.NewProfileRepository(db)
	userRepo := sharedRepo.NewUserRepository(db)

	// Service
	profileService := service.NewProfileService(profileRepo, userRepo, cacheRedis)

	// Handler
	profileHandler := handler.NewProfileHandler(profileService)

	// Setup handler to send to routes setup
	handler := &handler.Handlers{
		ProfileHandler: profileHandler,
	}

	routes.Setup(app, handler)
}

func InitApp() *sharedBootstrap.Deps {
	deps, err := sharedBootstrap.NewDeps("../../.env")
	if err != nil {
		log.Fatalf("Failed to connect to depedency: %v", err)
	}

	worker.StartLogWorker()

	helpers.InitLogger()

	middleware.InitWhitelistIP()

	return deps
}
