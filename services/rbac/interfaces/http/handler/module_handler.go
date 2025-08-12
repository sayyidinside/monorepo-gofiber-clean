package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type ModuleHandler interface {
	GetModule(c *fiber.Ctx) error
	GetAllModule(c *fiber.Ctx) error
	CreateModule(c *fiber.Ctx) error
	UpdateModule(c *fiber.Ctx) error
	DeleteModule(c *fiber.Ctx) error
}

type moduleHandler struct {
	service service.ModuleService
}

func NewModuleHandler(service service.ModuleService) ModuleHandler {
	return &moduleHandler{
		service: service,
	}
}

func (h *moduleHandler) GetModule(c *fiber.Ctx) error {
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

func (h *moduleHandler) GetAllModule(c *fiber.Ctx) error {
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

func (h *moduleHandler) CreateModule(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var input model.ModuleInput
	var response helpers.BaseResponse

	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		model.SanitizeModuleInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Errors:  err,
				Log:     &logData,
			})
		}

		response = h.service.Create(ctx, &input)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *moduleHandler) UpdateModule(c *fiber.Ctx) error {
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
		var input model.ModuleInput

		if err := c.BodyParser(&input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			})
		}

		model.SanitizeModuleInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
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

	return helpers.ResponseFormatter(c, response)
}

func (h *moduleHandler) DeleteModule(c *fiber.Ctx) error {
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
