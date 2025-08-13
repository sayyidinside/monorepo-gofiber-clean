package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Profile, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Profile, error)
	FindByUserID(ctx context.Context, user_id uint) (*entity.Profile, error)
	Insert(ctx context.Context, Profile *entity.Profile) error
	Update(ctx context.Context, Profile *entity.Profile) error
}

type profileRepository struct {
	*gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{DB: db}
}

func (r *profileRepository) FindByID(ctx context.Context, id uint) (*entity.Profile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var Profile entity.Profile
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Find(&Profile); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &Profile, nil
}

func (r *profileRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Profile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var Profile entity.Profile
	if err := r.DB.Limit(1).Where("uuid = ?", uuid).
		Find(&Profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &Profile, nil
}

func (r *profileRepository) FindByUserID(ctx context.Context, user_id uint) (*entity.Profile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var Profile entity.Profile
	if result := r.DB.WithContext(ctx).Limit(1).Where("user_id = ?", user_id).
		Find(&Profile); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &Profile, nil
}

func (r *profileRepository) Insert(ctx context.Context, Profile *entity.Profile) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(Profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *profileRepository) Update(ctx context.Context, Profile *entity.Profile) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", Profile.ID).Updates(Profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}
