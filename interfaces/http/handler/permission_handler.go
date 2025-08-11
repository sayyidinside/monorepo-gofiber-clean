package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type PermissionHandler interface {
	GetPermission(c *fiber.Ctx) error
	GetAllPermission(c *fiber.Ctx) error
	CreatePermission(c *fiber.Ctx) error
	UpdatePermission(c *fiber.Ctx) error
	DeletePermission(c *fiber.Ctx) error
}

type permissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) PermissionHandler {
	return &permissionHandler{
		service: service,
	}
}

func (h *permissionHandler) GetPermission(c *fiber.Ctx) error {
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

func (h *permissionHandler) GetAllPermission(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	query := new(model.QueryGet)
	if err := c.QueryParser(query); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request query",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		model.SanitizeQueryGet(query)

		url := c.BaseURL() + c.OriginalURL()
		response = h.service.GetAll(ctx, query, url)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *permissionHandler) CreatePermission(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	var input model.PermissionInput
	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		model.SanitizePermissionInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			return helpers.ResponseFormatter(c, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Errors:  err,
				Log:     &logData,
			})
		} else {
			response = h.service.Create(ctx, &input)
			response.Log = &logData
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *permissionHandler) UpdatePermission(c *fiber.Ctx) error {
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
		var input model.PermissionInput

		if err := c.BodyParser(&input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			})
		} else {
			model.SanitizePermissionInput(&input)

			if err := helpers.ValidateInput(input); err != nil {
				return helpers.ResponseFormatter(c, helpers.BaseResponse{
					Status:  fiber.StatusBadRequest,
					Success: false,
					Message: "Invalid or malformed request body",
					Errors:  err,
					Log:     &logData,
				})
			} else {
				response = h.service.UpdateByID(ctx, &input, uint(id))
				response.Log = &logData
			}
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *permissionHandler) DeletePermission(c *fiber.Ctx) error {
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
		response = h.service.DeleteByID(ctx, uint(id))
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
