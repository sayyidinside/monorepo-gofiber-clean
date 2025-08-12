package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	FindAllByUserID(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
	Insert(ctx context.Context, token *entity.RefreshToken) error
	RevokeByToken(ctx context.Context, token string) error
	RevokeAllByUserID(ctx context.Context, userID uint) error
	CountTokensByUserID(ctx context.Context, userID uint) (int64, error)
	DeleteExpiredTokens(ctx context.Context) error
}

type refreshTokenRepository struct {
	*gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{DB: db}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken

	result := r.DB.WithContext(ctx).Limit(1).Where("token = ?", token).Find(&refreshToken)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("refresh token not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindAllByUserID(ctx context.Context, userID uint) ([]entity.RefreshToken, error) {
	var tokens []entity.RefreshToken

	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (r *refreshTokenRepository) Insert(ctx context.Context, token *entity.RefreshToken) error {
	return r.DB.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	return r.DB.WithContext(ctx).Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uint) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error
}

func (r *refreshTokenRepository) CountTokensByUserID(ctx context.Context, userID uint) (int64, error) {
	var total int64

	if err := r.DB.WithContext(ctx).Model(&entity.RefreshToken{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *refreshTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.DB.WithContext(ctx).Where("expired_at < ?", time.Now()).Delete(&entity.RefreshToken{}).Error
}
