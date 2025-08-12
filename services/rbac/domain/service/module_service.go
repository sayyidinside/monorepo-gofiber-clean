package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type ModuleService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.ModuleInput) helpers.BaseResponse
	UpdateByID(ctx context.Context, input *model.ModuleInput, id uint) helpers.BaseResponse
	DeleteByID(ctx context.Context, id uint) helpers.BaseResponse
}

type moduleService struct {
	repository repository.ModuleRepository
}

func NewModuleService(repository repository.ModuleRepository) ModuleService {
	return &moduleService{
		repository: repository,
	}
}

func (s *moduleService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	module, err := s.repository.FindByID(ctx, id)
	if module == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Module not found",
		})
	}

	moduleModel := model.ModuleToDetailModel(module)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Module data found",
		Data:    moduleModel,
	}
}

func (s *moduleService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	modules, err := s.repository.FindAll(ctx, query)
	if modules == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Module not found",
			Errors:  err,
		})
	}

	moduleModels := model.ModuleToListModels(modules)

	totalData := s.repository.Count(ctx, query)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Module data found",
		Data:    moduleModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *moduleService) Create(ctx context.Context, input *model.ModuleInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	moduleEntity := model.ModuleInputToEntity(input)

	if err := s.validateEntityInput(ctx, moduleEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if err := s.repository.Insert(ctx, moduleEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "Module successfully created",
	})
}

func (s *moduleService) UpdateByID(ctx context.Context, input *model.ModuleInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check modul existence
	if module, err := s.repository.FindByID(ctx, id); module == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Module not found",
			Errors:  err,
		})
	}

	moduleEntity := model.ModuleInputToEntity(input)
	if moduleEntity == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		})
	}
	moduleEntity.ID = id

	if err := s.validateEntityInput(ctx, moduleEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if err := s.repository.Update(ctx, moduleEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Module succeessfully updated",
	})
}

func (s *moduleService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check modul existence
	module, err := s.repository.FindByID(ctx, id)
	if module == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Module not found",
			Errors:  err,
		})
	}

	if err := s.repository.Delete(ctx, module); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error deleting data",
			Errors:  err,
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Module successfully deleted",
	})
}

func (s *moduleService) validateEntityInput(ctx context.Context, module *entity.Module) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errs := []helpers.ValidationError{}

	// Check name duplication
	if exist := s.repository.NameExist(ctx, module); exist {
		errs = append(errs, helpers.ValidationError{
			Field: "name",
			Tag:   "duplicate",
		})
	}

	if len(errs) != 0 {
		logData.Message = "Validation error"
		logData.Err = errs
		return errs
	}

	return nil
}
