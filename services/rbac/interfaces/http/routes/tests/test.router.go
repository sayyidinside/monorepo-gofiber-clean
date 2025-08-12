package tests

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

func SetupApiTestRoutes(test fiber.Router) {
	test.Get("/panic", func(c *fiber.Ctx) error {
		panic("This panic is caught by fiber")
	})

	test.Get("/error", func(c *fiber.Ctx) error {
		return fmt.Errorf("Test error")
	})

	test.Get("/success", func(c *fiber.Ctx) error {
		// Sample data (you can pull this from a DB in a real scenario)

		// Mocked user data
		type user struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		users := []user{
			{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
			{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
		}

		data := interface{}(users)

		// Prepare pagination and meta information
		pagination := helpers.Pagination{
			CurrentPage: 1,
			TotalItems:  2,
			TotalPages:  1,
			ItemPerPage: 10,
			Self:        "/api/v1/users?page=1",
		}

		// Generating next and previous links
		next := "/api/v1/users?page=2"
		prev := (*string)(nil) // No previous page

		pagination.Next = &next
		pagination.Prev = prev

		// Prepare the response
		response := helpers.BaseResponse{
			Status:  200,
			Success: true,
			Message: "Data retrieved successfully",
			Data:    &data,
			Meta: &helpers.Meta{
				Pagination: &pagination,
			},
		}

		return helpers.ResponseFormatter(c, response)
	})

}
