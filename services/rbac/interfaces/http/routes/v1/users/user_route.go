package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/http/middleware"
)

func RegisterUserRoutes(route fiber.Router, handler handler.UserHandler) {
	user := route.Group("/data")

	user.Use(middleware.Authentication())

	user.Get(
		"/:id",
		middleware.Authorization(false, false, []string{
			"View User",
			"Create User",
			"Update User",
			"Delete User",
		}),
		handler.GetUser,
	)

	user.Get(
		"/",
		middleware.Authorization(false, false, []string{
			"View User",
			"Create User",
			"Update User",
			"Delete User",
		}),
		handler.GetAllUser,
	)

	user.Post(
		"/",
		middleware.Authorization(false, false, []string{
			"Create User",
		}),
		handler.CreateUser,
	)

	user.Put(
		"/:id/reset-password",
		middleware.Authorization(true, false, []string{}),
		handler.ResetPassword,
	)

	user.Put(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Update User",
		}),
		handler.UpdateUser,
	)

	user.Delete(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Delete User",
		}),
		handler.DeleteUser,
	)
}
