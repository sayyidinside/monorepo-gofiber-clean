package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type RoleService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.RoleInput) helpers.BaseResponse
	UpdateByID(ctx context.Context, input *model.RoleInput, id uint) helpers.BaseResponse
	DeleteByID(ctx context.Context, id uint) helpers.BaseResponse
}

type roleService struct {
	repository     repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(repository repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{
		repository:     repository,
		permissionRepo: permissionRepo,
	}
}

func (s *roleService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
			Errors:  err,
		})
	}

	roleModel := model.RoleToDetailModel(role)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role data found",
		Data:    roleModel,
	})
}

func (s *roleService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	roles, err := s.repository.FindAll(ctx, query)
	if roles == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
			Errors:  err,
		})
	}

	roleModels := model.RoleToListModels(roles)

	totalData := s.repository.Count(ctx, query)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role data found",
		Data:    roleModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *roleService) Create(ctx context.Context, input *model.RoleInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	roleEntity := model.RoleInputToEntity(input)

	if err := s.validateEntityInput(ctx, roleEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	permissions, err := s.permissionRepo.FindInID(ctx, input.Permissions)
	if err != nil || len(*permissions) == 0 {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Permission data not found",
			Errors:  err,
		})
	}

	roleEntity.Permissions = *permissions

	if err := s.repository.Insert(ctx, roleEntity); err != nil {
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
		Message: "Role successfully created",
	})
}

func (s *roleService) UpdateByID(ctx context.Context, input *model.RoleInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Start a new transaction
	tx := s.repository.BeginTransaction(ctx)

	// Check role existence
	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		tx.Rollback() // Rollback on error
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
			Errors:  err,
		})
	}

	roleEntity := model.RoleInputToEntity(input)
	if roleEntity == nil {
		tx.Rollback() // Rollback on error
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
			Errors:  err,
		})
	}

	roleEntity.ID = id

	// Retrieve permissions
	permissions, err := s.permissionRepo.FindInID(ctx, input.Permissions)
	if err != nil || len(*permissions) == 0 {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Permission data not found",
			Errors:  err,
		})
	}

	// Validate the entity
	if err := s.validateEntityInput(ctx, roleEntity); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	// Update role entity in the database
	if err := s.repository.UpdateWithTransaction(ctx, tx, roleEntity); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
			Errors:  err,
		})
	}

	// Replace permissions in the database
	if err := s.repository.ReplacePermissionsWithTransaction(ctx, tx, roleEntity, permissions); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error replacing role permissions data",
			Errors:  err,
		})
	}

	// Commit the transaction if all operations succeed
	tx.Commit()

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role successfully updated",
	})
}

func (s *roleService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check modul existence
	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
			Errors:  err,
		})
	}

	if err := s.repository.Delete(ctx, role); err != nil {
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
		Message: "Role successfully deleted",
	})
}

func (s *roleService) validateEntityInput(ctx context.Context, role *entity.Role) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errs := []helpers.ValidationError{}

	// Check name duplication
	if exist := s.repository.NameExist(ctx, role); exist {
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
