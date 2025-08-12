package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/shared/config"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

// Authentication to check and validate user access token, then set value from claim
// to context local as payload to be used in other layers
func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.AppConfig
		authorization := c.Get("Authorization")
		if authorization == "" {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Authorization token not provided",
			})
		}

		access_token := strings.TrimPrefix(authorization, "Bearer ")
		claim, err := helpers.ValidateToken(access_token, cfg.JwtAccessPublicSecret)
		if err != nil {
			if err.Error() == "validate: token has invalid claims: token is expired" {
				return helpers.ResponseFormatter(c, helpers.BaseResponse{
					Status:  fiber.StatusUnauthorized,
					Success: false,
					Message: "Token has expired",
				})
			}

			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		user_id, ok := claim["sub"].(float64)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		username, ok := claim["username"].(string)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		email, ok := claim["email"].(string)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		is_admin, ok := claim["is_admin"].(bool)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		validated, ok := claim["validated"].(bool)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		validated_at, ok := claim["validated_at"].(float64)
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		permissionInterfaces, ok := claim["permissions"].([]interface{})
		if !ok {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusUnauthorized,
				Success: false,
				Message: "Invalid token",
			})
		}

		var permissions []string
		for _, permission := range permissionInterfaces {
			if perm, ok := permission.(string); ok {
				permissions = append(permissions, perm)
			}
		}

		c.Locals("user_id", user_id)
		c.Locals("username", username)
		c.Locals("email", email)
		c.Locals("is_admin", is_admin)
		c.Locals("validated", validated)
		c.Locals("validated_at", time.Unix(int64(validated_at), 0))
		c.Locals("permissions", permissions)

		return c.Next()
	}
}

// Authorization middleware used to validate authenticated user have a permission to access endpoint.
//
// ! Important, that this middleware be called or used after Authentication middleware
func Authorization(isAdminOnly bool, isValidOnly bool, allowedPermissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// When option is admin true check user actually admin
		if isAdminOnly {
			if is_admin := c.Locals("is_admin").(bool); is_admin {
				return c.Next()
			}
		}

		// When option is valid true check user already validated
		if isValidOnly {
			validated := c.Locals("validated").(bool)
			validated_at := c.Locals("validated_at").(time.Time)

			if !validated || validated_at.After(time.Now()) {
				return helpers.ResponseFormatter(c, helpers.BaseResponse{
					Status:  fiber.StatusForbidden,
					Success: false,
					Message: "Unauthorized to access this resource",
				})
			}
		}

		// Check allowed permission against user permissions
		userPermissions := c.Locals("permissions").([]string)
		if len(userPermissions) > len(allowedPermissions) {
			// Create map from the smallest slice
			permissionMap := make(map[string]struct{}, len(allowedPermissions))
			for _, perm := range allowedPermissions {
				permissionMap[perm] = struct{}{}
			}

			// Check permission exists in the set
			for _, perm := range userPermissions {
				if _, exists := permissionMap[perm]; exists {
					return c.Next()
				}
			}
		} else {
			// Create map from the smallest slice
			permissionMap := make(map[string]struct{}, len(userPermissions))
			for _, perm := range userPermissions {
				permissionMap[perm] = struct{}{}
			}

			// Check permission exists in the set
			for _, perm := range allowedPermissions {
				if _, exists := permissionMap[perm]; exists {
					return c.Next()
				}
			}
		}

		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this resource",
		})
	}
}
