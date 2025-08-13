package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/interfaces/http/handler"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/interfaces/http/routes/v1/profile"
)

func RegisterRoutes(route fiber.Router, handler *handler.Handlers) {
	v1 := route.Group("/v1")

	profile.RegisterRoutes(v1, handler.ProfileHandler)
}
