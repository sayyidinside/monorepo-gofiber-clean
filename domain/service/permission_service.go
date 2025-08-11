package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type PermissionService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.PermissionInput) helpers.BaseResponse
	UpdateByID(ctx context.Context, input *model.PermissionInput, id uint) helpers.BaseResponse
	DeleteByID(ctx context.Context, id uint) helpers.BaseResponse
}

type permissionService struct {
	repository       repository.PermissionRepository
	moduleRepository repository.ModuleRepository
}

func NewPermissionService(repository repository.PermissionRepository, moduleRepository repository.ModuleRepository) PermissionService {
	return &permissionService{
		repository:       repository,
		moduleRepository: moduleRepository,
	}
}

func (s *permissionService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	permission, err := s.repository.FindByID(ctx, id)
	if permission == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
			Errors:  err,
		})
	}

	// convert entity to model data
	permissionModel := model.PermissionToDetailModel(permission)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Permission data found",
		Data:    permissionModel,
	})
}

func (s *permissionService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	permissions, err := s.repository.FindAll(ctx, query)
	if permissions == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
			Errors:  err,
		})
	}

	permissionModels := model.PermissionToListModels(permissions)

	totalData := s.repository.Count(ctx, query)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Permission data found",
		Data:    permissionModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *permissionService) Create(ctx context.Context, input *model.PermissionInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	permissionEntity := input.ToEntity()
	if permissionEntity == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		})
	}

	if err := s.validateEntityInput(ctx, permissionEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if err := s.repository.Insert(ctx, permissionEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
			Errors:  err,
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "Permission successfully created",
	})
}

func (s *permissionService) UpdateByID(ctx context.Context, input *model.PermissionInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check existence of permission
	if permission, err := s.repository.FindByID(ctx, id); permission == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
			Errors:  err,
		})
	}

	permissionEntity := input.ToEntity()
	if permissionEntity == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		})
	}
	permissionEntity.ID = id

	if err := s.validateEntityInput(ctx, permissionEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if err := s.repository.Update(ctx, permissionEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
			Errors:  err,
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Permission successfully updated",
	})
}

func (s *permissionService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check existence of permission
	permission, err := s.repository.FindByID(ctx, id)
	if permission == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
			Errors:  err,
		})
	}

	if err := s.repository.Delete(ctx, permission); err != nil {
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
		Message: "Permission successfully deleted",
	})
}

func (s *permissionService) validateEntityInput(ctx context.Context, permission *entity.Permission) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errs := []helpers.ValidationError{}

	// Check existence of module_id
	if module, err := s.moduleRepository.FindByID(ctx, permission.ModuleID); module == nil || err != nil {
		errs = append(errs, helpers.ValidationError{
			Field: "module_id",
			Tag:   "not_found",
		})
	}

	// Check name duplication
	if exist := s.repository.NameExist(ctx, permission); exist {
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
