package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	AllToken struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}

	RefreshTokenData struct {
		UUID      uuid.UUID `json:"uuid"`
		UserID    uint      `json:"user_id"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)
