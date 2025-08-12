package model

import (
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
)

type (
	RoleDetail struct {
		ID          uint              `json:"id"`
		UUID        uuid.UUID         `json:"uuid"`
		Name        string            `json:"name"`
		IsAdmin     bool              `json:"is_admin"`
		Permissions *[]PermissionList `json:"permissions"`
	}

	RoleList struct {
		ID   uint      `json:"id"`
		UUID uuid.UUID `json:"uuid"`
		Name string    `json:"name"`
	}

	RoleInput struct {
		Name        string `json:"name" form:"name" xml:"name" validate:"required"`
		IsAdmin     bool   `json:"is_admin" form:"is_admin" xml:"is_admin" validate:"boolean"`
		Permissions []uint `json:"permissions" form:"permissions" xml:"permissions" validate:"required,gt=0,dive,numeric"`
	}
)

func RoleToDetailModel(role *entity.Role) *RoleDetail {
	permissions := PermissionToListModels(&role.Permissions)

	return &RoleDetail{
		ID:          role.ID,
		UUID:        role.UUID,
		Name:        role.Name,
		Permissions: permissions,
	}
}

func RoleToListModel(role *entity.Role) *RoleList {
	return &RoleList{
		ID:   role.ID,
		UUID: role.UUID,
		Name: role.Name,
	}
}

func RoleToListModels(role *[]entity.Role) *[]RoleList {
	listModels := []RoleList{}

	for _, role := range *role {
		listModels = append(listModels, *RoleToListModel(&role))
	}

	return &listModels
}

func RoleInputToEntity(input *RoleInput) *entity.Role {
	return &entity.Role{
		Name:    input.Name,
		IsAdmin: input.IsAdmin,
	}
}

func SanitizeRoleInput(input *RoleInput) {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
}
