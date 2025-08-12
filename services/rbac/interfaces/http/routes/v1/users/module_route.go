package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/middleware"
)

func RegisterModuleRoutes(route fiber.Router, handler handler.ModuleHandler) {
	modules := route.Group("/modules")

	modules.Use(middleware.Authentication())

	modules.Get(
		"/",
		middleware.Authorization(true, true, []string{}),
		handler.GetAllModule,
	)

	modules.Get(
		"/:id",
		middleware.Authorization(true, true, []string{}),
		handler.GetModule,
	)

	modules.Post(
		"",
		middleware.Authorization(true, true, []string{}),
		handler.CreateModule,
	)

	modules.Put(
		"/:id",
		middleware.Authorization(true, true, []string{}),
		handler.UpdateModule,
	)

	modules.Delete(
		"/:id",
		middleware.Authorization(true, true, []string{}),
		handler.DeleteModule,
	)
}
