package model

import "github.com/microcosm-cc/bluemonday"

type (
	LoginInput struct {
		UsernameOrEmail string `json:"username_or_email" form:"username_or_email" xml:"username_or_email" validate:"required"`
		Password        string `json:"password" form:"password" xml:"password" validate:"required"`
	}

	TokenInput struct {
		Token string `json:"token" form:"token" xml:"token" validate:"required"`
	}
)

func (input *LoginInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.UsernameOrEmail = sanitizer.Sanitize(input.UsernameOrEmail)
	input.Password = sanitizer.Sanitize(input.Password)
}

func (input *TokenInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Token = sanitizer.Sanitize(input.Token)
}
