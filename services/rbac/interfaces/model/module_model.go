package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
)

type (
	ModuleDetail struct {
		ID          uint              `json:"id"`
		UUID        uuid.UUID         `json:"uuid"`
		Name        string            `json:"name"`
		Permissions *[]PermissionList `json:"permissions"`
		CreatedAt   time.Time         `json:"created_at"`
		UpdatedAt   time.Time         `json:"updated_at"`
	}

	ModuleList struct {
		ID   uint      `json:"id"`
		UUID uuid.UUID `json:"uuid"`
		Name string    `json:"name"`
	}

	ModuleInput struct {
		Name string `json:"name" form:"name" xml:"name" validate:"required"`
	}
)

func ModuleToDetailModel(module *entity.Module) *ModuleDetail {
	permissions := PermissionToListModels(&module.Permissions)

	return &ModuleDetail{
		ID:          module.ID,
		UUID:        module.UUID,
		Name:        module.Name,
		Permissions: permissions,
	}
}

func ModuleToListModel(module *entity.Module) *ModuleList {
	return &ModuleList{
		ID:   module.ID,
		UUID: module.UUID,
		Name: module.Name,
	}
}

func ModuleToListModels(modules *[]entity.Module) *[]ModuleList {
	listModels := []ModuleList{}

	for _, module := range *modules {
		listModels = append(listModels, *ModuleToListModel(&module))
	}

	return &listModels
}

func SanitizeModuleInput(input *ModuleInput) {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
}

func ModuleInputToEntity(input *ModuleInput) *entity.Module {
	return &entity.Module{
		Name: input.Name,
	}
}
