package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
)

type (
	PermissionDetail struct {
		ID        uint      `json:"id"`
		UUID      uuid.UUID `json:"uuid"`
		Name      string    `json:"name"`
		Module    string    `json:"module"`
		ModuleID  uint      `json:"module_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PermissionList struct {
		ID       uint      `json:"id"`
		UUID     uuid.UUID `json:"uuid"`
		Name     string    `json:"name"`
		Module   string    `json:"module"`
		ModuleID uint      `json:"module_id"`
	}

	PermissionInput struct {
		Name     string `json:"name" form:"name" xml:"name" validate:"required"`
		ModuleID uint   `json:"module_id" form:"module_id" xml:"module_id" validate:"required,numeric"`
	}
)

func PermissionToDetailModel(permission *entity.Permission) *PermissionDetail {
	return &PermissionDetail{
		ID:        permission.ID,
		UUID:      permission.UUID,
		Name:      permission.Name,
		Module:    permission.Module.Name,
		ModuleID:  permission.Module.ID,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

func PermissionToListModel(permission *entity.Permission) *PermissionList {
	return &PermissionList{
		ID:       permission.ID,
		UUID:     permission.UUID,
		Name:     permission.Name,
		ModuleID: permission.ModuleID,
		Module:   permission.Module.Name,
	}
}

func PermissionToListModels(permissions *[]entity.Permission) *[]PermissionList {
	listModels := []PermissionList{}

	for _, permission := range *permissions {
		listModels = append(listModels, *PermissionToListModel(&permission))
	}

	return &listModels
}

func SanitizePermissionInput(input *PermissionInput) {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
}

func (input *PermissionInput) ToEntity() *entity.Permission {

	return &entity.Permission{
		Name:     input.Name,
		ModuleID: input.ModuleID,
	}
}
