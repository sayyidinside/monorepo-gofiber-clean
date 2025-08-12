package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/http/middleware"
)

func RegisterRoutes(route fiber.Router, handler handler.AuthHandler) {
	authRoutes := route.Group("/auth")

	authRoutes.Post("/login", handler.Login)
	authRoutes.Get("/verify", handler.Verify)
	authRoutes.Post("/refresh", handler.Refresh)
	authRoutes.Post("/logout", middleware.Authentication(), handler.Logout)
}
