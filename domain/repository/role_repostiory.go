package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type RoleRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByID(ctx context.Context, id uint) (*entity.Role, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*entity.Role, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Role, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Role, error)
	Count(ctx context.Context, query *model.QueryGet) int64
	CountUnscoped(ctx context.Context, query *model.QueryGet) int64
	Insert(ctx context.Context, role *entity.Role) error
	UpdateWithTransaction(ctx context.Context, tx *gorm.DB, role *entity.Role) error
	Delete(ctx context.Context, role *entity.Role) error
	NameExist(ctx context.Context, role *entity.Role) bool
	ReplacePermissionsWithTransaction(ctx context.Context, tx *gorm.DB, role *entity.Role, permissions *[]entity.Permission) error
}

type roleRepository struct {
	*gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.DB.Begin()
}

func (r *roleRepository) FindByID(ctx context.Context, id uint) (*entity.Role, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var role entity.Role
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "uuid", "module_id")
		}).
		Preload("Permissions.Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&role); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &role, nil
}

func (r *roleRepository) FindByIDUnscoped(ctx context.Context, id uint) (*entity.Role, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var role entity.Role
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).Unscoped().
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "uuid", "module_id")
		}).
		Preload("Permissions.Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&role); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &role, nil
}

func (r *roleRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Role, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var role entity.Role

	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "uuid", "module_id")
		}).
		Preload("Permissions.Module", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&role); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error

		return nil, result.Error
	}

	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Role, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var roles []entity.Role

	tx := r.DB.WithContext(ctx).Model(&entity.Role{})

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":     "name",
		"is_admin": "is_admin",
		"updated":  "updated_at",
		"created":  "created_at",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&roles).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &roles, nil
}

func (r *roleRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Role{})

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":     "name",
		"is_admin": "is_admin",
		"updated":  "updated_at",
		"created":  "created_at",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	tx.Count(&total)

	return total
}

func (r *roleRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Role{}).Unscoped()

	// map value for parsing user query input
	var allowedFields = map[string]string{
		"name":     "name",
		"is_admin": "is_admin",
		"updated":  "updated_at",
		"created":  "created_at",
	}

	// Apply Query Operation
	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	tx.Count(&total)

	return total
}

func (r *roleRepository) Insert(ctx context.Context, role *entity.Role) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return r.DB.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", role.ID).Updates(role).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *roleRepository) UpdateWithTransaction(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", role.ID).Updates(role).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, role *entity.Role) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", role.ID).Delete(role).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *roleRepository) NameExist(ctx context.Context, role *entity.Role) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Role{}).Where("name = ?", role.Name)

	if role.ID != 0 {
		tx = tx.Not("id = ?", role.ID)
	}

	if result := tx.Count(&total); result.Error != nil {
		logData.Err = result.Error
		logData.Message = "Not Passed"
	}

	return total != 0
}

func (r *roleRepository) ReplacePermissionsWithTransaction(ctx context.Context, tx *gorm.DB, role *entity.Role, permissions *[]entity.Permission) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// GORM requires data to be pre-loaded before using the Association.
	tx.Preload("Permissions").First(&role)

	tx.Model(&role).Association("Permissions").Replace(permissions)
	if tx.Error != nil {
		logData.Message = "Not Passed"
		logData.Err = tx.Error
		return tx.Error
	}

	return nil

}
