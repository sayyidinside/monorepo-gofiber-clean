package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
)

type (
	ProfileDetail struct {
		ID         uint      `json:"id"`
		UUID       uuid.UUID `json:"uuid"`
		UserID     uint      `json:"user_id"`
		Name       string    `json:"name"`
		Email      string    `json:"email"`
		BirthPlace string    `json:"birth_place"`
		BirthDate  time.Time `json:"birth_date"`
		Address    string    `json:"address"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	ProfileInput struct {
		Name       string `json:"name" form:"name" xml:"name" validate:"required"`
		BirthPlace string `json:"birth_place" form:"birth_place" xml:"birth_place" validate:"required"`
		BirthDate  string `json:"birth_date" form:"birth_date" xml:"birth_date" validate:"required,datetime=2006-01-02"`
		Address    string `json:"address" form:"address" xml:"address" validate:"required"`
	}

	ProfileInputByAdmin struct {
		UserID     uint   `json:"user_id" form:"user_id" xml:"user_id" validate:"required"`
		Name       string `json:"name" form:"name" xml:"name" validate:"required"`
		BirthPlace string `json:"birth_place" form:"birth_place" xml:"birth_place" validate:"required"`
		BirthDate  string `json:"birth_date" form:"birth_date" xml:"birth_date" validate:"required,datetime=2006-01-02"`
		Address    string `json:"address" form:"address" xml:"address" validate:"required"`
	}
)

func ProfileToDetailModel(profile *entity.Profile) *ProfileDetail {
	return &ProfileDetail{
		ID:         profile.ID,
		UUID:       profile.UUID,
		UserID:     profile.UserID,
		Name:       profile.Name,
		Email:      profile.User.Email,
		BirthPlace: profile.BirthPlace,
		BirthDate:  profile.BirthDate,
		Address:    profile.Address,
		CreatedAt:  profile.CreatedAt,
		UpdatedAt:  profile.UpdatedAt,
	}
}

func (input *ProfileInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
	input.BirthPlace = sanitizer.Sanitize(input.BirthPlace)
	input.BirthDate = sanitizer.Sanitize(input.BirthDate)
	input.Address = sanitizer.Sanitize(input.Address)
}

func (input *ProfileInputByAdmin) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
	input.BirthPlace = sanitizer.Sanitize(input.BirthPlace)
	input.BirthDate = sanitizer.Sanitize(input.BirthDate)
	input.Address = sanitizer.Sanitize(input.Address)
}

func (input *ProfileInput) ToEntity() (*entity.Profile, error) {
	birth_date, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Name:       input.Name,
		BirthPlace: input.BirthPlace,
		BirthDate:  birth_date,
		Address:    input.Address,
	}, nil
}

func (input *ProfileInputByAdmin) ToEntity() (*entity.Profile, error) {
	birth_date, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		UserID:     input.UserID,
		Name:       input.Name,
		BirthPlace: input.BirthPlace,
		BirthDate:  birth_date,
		Address:    input.Address,
	}, nil
}
