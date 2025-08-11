package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"gorm.io/gorm"
)

type RefreshToken struct {
	UUID      uuid.UUID `json:"uuid" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token" gorm:"index,type:longtext"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field.
func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	cfg := config.AppConfig

	if r.UUID == uuid.Nil {
		r.UUID = uuid.New()
	}

	r.ExpiredAt = time.Now().Add(time.Duration(cfg.JwtRefreshTime) * time.Hour)
	return
}
