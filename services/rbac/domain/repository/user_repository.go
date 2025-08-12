package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.User, error)
	Count(ctx context.Context, query *model.QueryGet) int64
	CountUnscoped(ctx context.Context, query *model.QueryGet) int64
	Insert(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
	NameExist(ctx context.Context, user *entity.User) bool
	EmailExist(ctx context.Context, user *entity.User) bool
	UsernameExist(ctx context.Context, user *entity.User) bool
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error)
	// Create(*User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var user entity.User
	result := r.DB.WithContext(ctx).
		Limit(1).
		Where("id = ?", id).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user)

	if result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var user entity.User
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.User, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var users []entity.User
	tx := r.DB.WithContext(ctx).Model(&entity.User{}).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		})

	var allowedFields = map[string]string{
		"name":    "users.name",
		"created": "users.created_at",
		"updated": "users.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&users).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &users, nil
}

func (r *userRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		})

	var allowedFields = map[string]string{
		"name":    "users.name",
		"created": "users.created_at",
		"updated": "users.updated_at",
	}

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

func (r *userRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		})

	var allowedFields = map[string]string{
		"name":    "users.name",
		"created": "users.created_at",
		"updated": "users.updated_at",
	}

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

func (r *userRepository) Insert(ctx context.Context, user *entity.User) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, user *entity.User) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(user).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *userRepository) NameExist(ctx context.Context, user *entity.User) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).Where("username = ? ", user.Username)
	if user.ID != 0 {
		tx = tx.Not("id = ?", user.ID)
	}

	if err := tx.Count(&totalData).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
	}

	return totalData != 0
}

func (r *userRepository) EmailExist(ctx context.Context, user *entity.User) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).Where("email = ?", user.Email)
	if user.ID != 0 {
		tx = tx.Not("id = ?", user.ID)
	}

	if err := tx.Count(&totalData).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
	}
	return totalData != 0
}

func (r *userRepository) UsernameExist(ctx context.Context, user *entity.User) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).Where("username = ?", user.Username)
	if user.ID != 0 {
		tx = tx.Not("id = ?", user.ID)
	}

	tx.Count(&totalData)
	if err := tx.Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
	}

	return totalData != 0
}

func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*entity.User, error) {
	var user entity.User

	result := r.DB.WithContext(ctx).Limit(1).Where("username = ?", usernameOrEmail).Or("email = ?", usernameOrEmail).Preload("Role").Preload("Role.Permissions").Find(&user)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user data not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
