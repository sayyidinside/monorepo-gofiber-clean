package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/middleware"
)

func RegisterRoleRoutes(route fiber.Router, handler handler.RoleHandler) {
	role := route.Group("/roles")

	role.Use(middleware.Authentication())

	role.Get(
		"/:id",
		middleware.Authorization(false, false, []string{
			"View Role",
			"Create Role",
			"Update Role",
			"Delete Role",
		}),
		handler.GetRole,
	)

	role.Get(
		"/",
		middleware.Authorization(false, false, []string{
			"View Role",
			"Create Role",
			"Update Role",
			"Delete Role",
		}),
		handler.GetAllRole,
	)

	role.Post(
		"",
		middleware.Authorization(false, false, []string{
			"Create Role",
		}),
		handler.CreateRole,
	)

	role.Put(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Update Role",
		}),
		handler.UpdateRole,
	)

	role.Delete(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Delete Role",
		}),
		handler.DeleteRole,
	)
}
