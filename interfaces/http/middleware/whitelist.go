package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

// Global variable to hold allowed IPs set
var allowedIPSet map[string]struct{}

// init function to initialize allowedIPSet once
func InitWhitelistIP() {
	cfg := config.AppConfig
	allowedIPs := strings.Split(cfg.AllowedIPs, ",")
	allowedIPSet = make(map[string]struct{}, len(allowedIPs))

	// Populate the map with allowed IPs
	for _, ip := range allowedIPs {
		allowedIPSet[ip] = struct{}{}
	}

	// Add localhost to the allowed IPs
	allowedIPSet["127.0.0.1"] = struct{}{}
}

// WhitelistIP is the middleware handler to check IPs
func WhitelistIP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Jika semua IP diizinkan, langsung lanjut ke handler berikutnya
		if config.AppConfig.AllowedIPs == "*" {
			return c.Next()
		}

		clientIP := c.IP()

		// Cek apakah IP ada dalam set allowedIPSet
		if _, ok := allowedIPSet[clientIP]; ok {
			return c.Next() // IP diizinkan
		}

		// Kembalikan response unauthorized jika IP tidak diizinkan
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  401,
			Success: false,
			Message: "Unauthorized",
		})
	}
}
