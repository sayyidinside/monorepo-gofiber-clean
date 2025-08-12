package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Permission, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Permission, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Permission, error)
	FindInID(ctx context.Context, ids []uint) (*[]entity.Permission, error)
	Insert(ctx context.Context, permission *entity.Permission) error
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, permission *entity.Permission) error
	Count(ctx context.Context, query *model.QueryGet) int64
	CountUnscoped(ctx context.Context, query *model.QueryGet) int64
	NameExist(ctx context.Context, permission *entity.Permission) bool
}

type permissionRepository struct {
	*gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{DB: db}
}

func (r *permissionRepository) FindByID(ctx context.Context, id uint) (*entity.Permission, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var permission entity.Permission
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Preload("Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&permission); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &permission, nil
}

func (r *permissionRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Permission, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var permission entity.Permission
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&permission); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &permission, nil
}

func (r *permissionRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Permission, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var permissions []entity.Permission

	tx := r.DB.WithContext(ctx).Model(&entity.Permission{}).
		Joins("JOIN modules on modules.id = permissions.module_id").
		Preload("Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		})

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":        "permissions.name",
		"module":      "permissions.module_id",
		"updated":     "permissions.updated_at",
		"created":     "permissions.created_at",
		"module_name": "modules.name",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&permissions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &permissions, nil
}

func (r *permissionRepository) FindInID(ctx context.Context, ids []uint) (*[]entity.Permission, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var permissions []entity.Permission

	if err := r.DB.WithContext(ctx).Model(&entity.Permission{}).Select("id", "name", "module_id").Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &permissions, nil
}

func (r *permissionRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Permission{}).
		Joins("JOIN modules on modules.id = permissions.module_id")

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":        "permissions.name",
		"module":      "permissions.module_id",
		"updated":     "permissions.updated_at",
		"created":     "permissions.created_at",
		"module_name": "modules.name",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return total
}

func (r *permissionRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Permission{}).Unscoped().
		Joins("JOIN modules on modules.id = permissions.module_id")

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":        "permissions.name",
		"module":      "permissions.module_id",
		"updated":     "permissions.updated_at",
		"created":     "permissions.created_at",
		"module_name": "modules.name",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return total
}

func (r *permissionRepository) Insert(ctx context.Context, permission *entity.Permission) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return r.DB.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", permission.ID).Updates(permission).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, permission *entity.Permission) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(permission).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *permissionRepository) NameExist(ctx context.Context, permission *entity.Permission) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var totalData int64
	tx := r.DB.WithContext(ctx).Model(&entity.Permission{}).Where("name = ? AND module_id = ?", permission.Name, permission.ModuleID)

	if permission.ID != 0 {
		tx = tx.Not("id = ?", permission.ID)
	}

	if err := tx.Count(&totalData).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
	}

	return totalData != 0
}
