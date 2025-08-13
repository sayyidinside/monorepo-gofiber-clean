package bootstrap

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/rbac/cmd/worker"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/rbac/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/rbac/interfaces/http/routes"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/service"
	sharedBootstrap "github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/bootstrap"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/http/middleware"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"gorm.io/gorm"
)

func Initialize(app *fiber.App, db *gorm.DB, cacheRedis *redis.CacheClient, lockRedis *redis.LockClient) {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Service
	userService := service.NewUserService(userRepo, roleRepo, cacheRedis)
	permissionService := service.NewPermissionService(permissionRepo, moduleRepo)
	moduleService := service.NewModuleService(moduleRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	authService := service.NewAuthService(refreshTokenRepo, userRepo)

	// Handler
	userHandler := handler.NewUserHandler(userService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	moduleHandler := handler.NewModuleHandler(moduleService)
	roleHandler := handler.NewRoleHandler(roleService)
	authHandler := handler.NewAuthHandler(authService)

	// Setup handler to send to routes setup
	handler := &handler.Handlers{
		UserManagementHandler: &handler.UserManagementHandler{
			UserHandler:       userHandler,
			PermissionHandler: permissionHandler,
			ModuleHandler:     moduleHandler,
			RoleHandler:       roleHandler,
		},
		AuthHandler: authHandler,
	}

	routes.Setup(app, handler)
}

func InitApp() *sharedBootstrap.Deps {
	deps, err := sharedBootstrap.NewDeps()
	if err != nil {
		log.Fatalf("Failed to connect to depedency: %v", err)
	}

	worker.StartLogWorker()

	helpers.InitLogger()

	middleware.InitWhitelistIP()

	return deps
}
