package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/pkg/helpers"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	logDataData := helpers.CreateLog(c)
	var input model.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logDataData,
		})
	}

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logDataData,
		})
	}

	input.Sanitize()

	response := h.service.Login(c.Context(), &input)
	response.Log = &logDataData

	return helpers.ResponseFormatter(c, response)
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	var input model.TokenInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	input.Sanitize()

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.Logout(c.Context(), input.Token)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *authHandler) Refresh(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	var input model.TokenInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	input.Sanitize()

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.Refresh(c.Context(), input.Token)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *authHandler) Verify(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	var input model.TokenInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	input.Sanitize()

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.VerifyAccessToken(c.Context(), input.Token)
	if !response.Success {
		response = h.service.VerifyRefreshToken(c.Context(), input.Token)
	}
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}
