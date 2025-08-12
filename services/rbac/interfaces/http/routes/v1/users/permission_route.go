package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/rbac/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/http/middleware"
)

func RegisterPermissionRoutes(route fiber.Router, handler handler.PermissionHandler) {
	permission := route.Group("/permissions")

	permission.Use(middleware.Authentication())

	permission.Get(
		"/",
		middleware.Authorization(false, false, []string{
			"View Permission",
			"Create Permission",
			"Update Permission",
			"Delete Permission",
		}),
		handler.GetAllPermission,
	)

	permission.Get(
		"/:id",
		middleware.Authorization(false, false, []string{
			"View Permission",
			"Create Permission",
			"Update Permission",
			"Delete Permission",
		}),
		handler.GetPermission,
	)

	permission.Post(
		"",
		middleware.Authorization(false, false, []string{
			"Create Permission",
		}),
		handler.CreatePermission,
	)

	permission.Put(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Update Permission",
		}),
		handler.UpdatePermission,
	)
	permission.Delete(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Delete Permission",
		}),
		handler.DeletePermission,
	)
}
