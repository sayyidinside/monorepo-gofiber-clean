package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Module, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*entity.Module, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Module, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Module, error)
	Count(ctx context.Context, query *model.QueryGet) int64
	CountUnscoped(ctx context.Context, query *model.QueryGet) int64
	Insert(ctx context.Context, module *entity.Module) error
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, module *entity.Module) error
	NameExist(ctx context.Context, module *entity.Module) bool
}

type moduleRepository struct {
	*gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{DB: db}
}

func (r *moduleRepository) FindByID(ctx context.Context, id uint) (*entity.Module, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var module entity.Module
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "uuid", "module_id")
		}).
		Find(&module); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &module, nil
}

func (r *moduleRepository) FindByIDUnscoped(ctx context.Context, id uint) (*entity.Module, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var module entity.Module
	if err := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).Unscoped().
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "module_id").Unscoped()
		}).
		Find(&module).Error; err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Module, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var module entity.Module
	if err := r.DB.Limit(1).Where("uuid = ?", uuid).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "module_id").Unscoped()
		}).
		Find(&module).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Module, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var modules []entity.Module

	tx := r.DB.WithContext(ctx).Model(&entity.Module{})

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":    "name",
		"updated": "updated_at",
		"created": "created_at",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&modules).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &modules, nil
}

func (r *moduleRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Module{})

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":    "name",
		"updated": "updated_at",
		"created": "created_at",
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

func (r *moduleRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Module{}).Unscoped()

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":    "name",
		"updated": "updated_at",
		"created": "created_at",
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

func (r *moduleRepository) Insert(ctx context.Context, module *entity.Module) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(module).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *moduleRepository) Update(ctx context.Context, module *entity.Module) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", module.ID).Updates(module).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *moduleRepository) Delete(ctx context.Context, module *entity.Module) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(module).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *moduleRepository) NameExist(ctx context.Context, module *entity.Module) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Module{}).Where("name = ?", module.Name)

	if module.ID != 0 {
		tx = tx.Not("id = ?", module.ID)
	}

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return total != 0
}
