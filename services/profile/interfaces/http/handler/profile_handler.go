package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type ProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	GetProfileByAdmin(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UpdateProfileByAdmin(c *fiber.Ctx) error
}

type profileHandler struct {
	service service.ProfileService
}

func NewProfileHandler(service service.ProfileService) ProfileHandler {
	return &profileHandler{
		service: service,
	}
}

func (h *profileHandler) GetProfile(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	user_id := c.Locals("user_id").(float64)
	response := h.service.GetByUserID(ctx, uint(user_id))
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *profileHandler) GetProfileByAdmin(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID Format",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		response = h.service.GetByUserUUID(ctx, uuid)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *profileHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse
	user_id := c.Locals("user_id").(float64)

	var input model.ProfileInput

	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	}

	input.Sanitize()

	if err := helpers.ValidateInput(input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	} else {
		response = h.service.UpdateByUserID(ctx, &input, uint(user_id))
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *profileHandler) UpdateProfileByAdmin(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid UUID Format",
			Log:     &logData,
			Errors:  err,
		})
	}

	var input model.ProfileInputByAdmin

	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	}

	input.Sanitize()

	if err := helpers.ValidateInput(input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	} else {
		response = h.service.UpdateByUserUUID(ctx, &input, uuid)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
