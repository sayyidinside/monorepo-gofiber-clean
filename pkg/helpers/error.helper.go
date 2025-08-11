package helpers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// NotFoundHelpers handles 404 - Route not found
func NotFoundHelper(c *fiber.Ctx) error {
	logData := CreateLog(c)
	return ResponseFormatter(c, BaseResponse{
		Status:  404,
		Success: false,
		Message: "Resource Not Found",
		Log:     &logData,
	})
}

// ErrorHelper handles unhandled errors (500)
func ErrorHelper(c *fiber.Ctx) error {
	logData := CreateLog(c)
	// Try to handle the request and capture any unhandled errors
	err := c.Next() // Process next middleware or route handler
	if err != nil {
		// Log the error for debugging
		log.Printf("Unhandled error: %v", err)

		// Return a 500 Internal Server Error response
		return ResponseFormatter(c, BaseResponse{
			Status:  500,
			Success: false,
			Message: "Internal Server Error",
			Errors:  err,
			Log:     &logData,
		})
	}
	return nil
}

func RecoverWithLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logData := CreateLog(c)
				// Log the panic details
				log.Printf("Panic: %v\n", r)

				// Send the panic error to ErrorHandler
				ResponseFormatter(c, BaseResponse{
					Status:  500,
					Success: false,
					Message: "Internal Server Error",
					Errors:  r,
					Log:     &logData,
				})
			}
		}()

		return c.Next()
	}
}
