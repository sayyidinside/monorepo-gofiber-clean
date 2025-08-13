package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type RoleHandler interface {
	GetRole(c *fiber.Ctx) error
	GetAllRole(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

type roleHandler struct {
	service service.RoleService
}

func NewRoleHandler(service service.RoleService) RoleHandler {
	return &roleHandler{
		service: service,
	}
}

func (h *roleHandler) GetRole(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		response = h.service.GetByID(ctx, uint(id))
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
func (h *roleHandler) GetAllRole(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	query := new(model.QueryGet)
	var response helpers.BaseResponse

	if err := c.QueryParser(query); err != nil {
		response = helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request query",
			Log:     &logData,
			Errors:  err,
		}
	} else {
		model.SanitizeQueryGet(query)
		url := c.BaseURL() + c.OriginalURL()
		response = h.service.GetAll(ctx, query, url)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *roleHandler) CreateRole(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var input model.RoleInput
	var response helpers.BaseResponse

	if err := c.BodyParser(&input); err != nil {
		response = helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		}
	} else {
		model.SanitizeRoleInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Errors:  err,
				Log:     &logData,
			}
		} else {
			response = h.service.Create(ctx, &input)
			response.Log = &logData
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *roleHandler) UpdateRole(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var response helpers.BaseResponse

	if err != nil {
		response = helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		}
	} else {
		var input model.RoleInput

		if err := c.BodyParser(&input); err != nil {
			response = helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			}
		} else {
			model.SanitizeRoleInput(&input)

			if err := helpers.ValidateInput(input); err != nil {
				response = helpers.BaseResponse{
					Status:  fiber.StatusBadRequest,
					Success: false,
					Message: "Invalid or malformed request body",
					Errors:  err,
					Log:     &logData,
				}
			} else {
				response = h.service.UpdateByID(ctx, &input, uint(id))
				response.Log = &logData
			}
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *roleHandler) DeleteRole(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var response helpers.BaseResponse

	if err != nil {
		response = helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		}
	} else {
		response = h.service.DeleteByID(ctx, uint(id))
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
