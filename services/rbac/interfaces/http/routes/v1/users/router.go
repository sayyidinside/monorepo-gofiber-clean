package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/handler"
)

func RegisterRoutes(route fiber.Router, handler *handler.UserManagementHandler) {
	user := route.Group("/users/")

	RegisterUserRoutes(user, handler.UserHandler)
	RegisterPermissionRoutes(user, handler.PermissionHandler)
	RegisterModuleRoutes(user, handler.ModuleHandler)
	RegisterRoleRoutes(user, handler.RoleHandler)
}
