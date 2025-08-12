package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/http/middleware"
)

func RegisterRoutes(route fiber.Router, handler handler.ProfileHandler) {
	profileRoutes := route.Group("/profile")

	profileRoutes.Use(middleware.Authentication())

	profileRoutes.Get("/", handler.GetProfile)
	profileRoutes.Get(
		"/:uuid",
		middleware.Authorization(true, true, []string{}),
		handler.GetProfileByAdmin,
	)
	profileRoutes.Put(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.UpdateProfile,
	)
	profileRoutes.Put(
		"/:uuid",
		middleware.Authorization(true, true, []string{}),
		handler.UpdateProfileByAdmin,
	)
}
