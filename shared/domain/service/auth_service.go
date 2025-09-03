package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, input *model.LoginInput) helpers.BaseResponse
	Logout(ctx context.Context, refreshToken string) helpers.BaseResponse
	Refresh(ctx context.Context, refreshToken string) helpers.BaseResponse
	VerifyAccessToken(ctx context.Context, accessToken string) helpers.BaseResponse
	VerifyRefreshToken(ctx context.Context, refreshToken string) helpers.BaseResponse
}

type authService struct {
	refreshTokenRepository repository.RefreshTokenRepository
	userRepository         repository.UserRepository
	broadcastService       BroadcastService
}

func NewAuthService(
	refreshTokenRepository repository.RefreshTokenRepository, userRepository repository.UserRepository,
	broadcastService BroadcastService,
) AuthService {
	return &authService{
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
		broadcastService:       broadcastService,
	}
}

func (s *authService) Login(ctx context.Context, input *model.LoginInput) helpers.BaseResponse {
	cfg := config.AppConfig

	user, err := s.userRepository.FindByUsernameOrEmail(ctx, input.UsernameOrEmail)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid email / username or password",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid email / username or password",
		}
	}

	accessToken, err := helpers.GenerateToken(user, cfg.JwtAccessTime, cfg.JwtAccessPrivateSecret, false)
	if err != nil {
		jsonData, _ := json.MarshalIndent(err.Error(), "", " ")
		log.Println(string(jsonData))
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	refreshToken, err := helpers.GenerateToken(user, cfg.JwtRefreshTime, cfg.JwtRefreshPrivateSecret, true)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	if err := s.refreshTokenRepository.Insert(ctx, &entity.RefreshToken{UserID: user.ID, Token: refreshToken}); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed registering token",
		}
	}

	emailContent, err := model.ToTemplateEmail("Login", "Success")
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed constructing email content",
		}
	}

	s.broadcastService.SendEmail(ctx, user.ID, "New Login Information", emailContent)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("%s user successfully login", user.Username),
		Data: &model.AllToken{
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		},
	}
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) helpers.BaseResponse {
	cfg := config.AppConfig

	_, err := helpers.ValidateToken(refreshToken, cfg.JwtRefreshPublicSecret)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	tokenEntity, err := s.refreshTokenRepository.FindByToken(ctx, refreshToken)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	user, err := s.userRepository.FindByID(ctx, tokenEntity.UserID)
	if err != nil || user == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "User not found",
		}
	}

	user, err = s.userRepository.FindByUsernameOrEmail(ctx, user.Username)
	if err != nil || user == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "User not found",
		}
	}

	newAccessToken, err := helpers.GenerateToken(user, cfg.JwtAccessTime, cfg.JwtAccessPrivateSecret, false)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	newRefreshToken, err := helpers.GenerateToken(user, cfg.JwtRefreshTime, cfg.JwtRefreshPrivateSecret, true)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	if err := s.refreshTokenRepository.RevokeByToken(ctx, refreshToken); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed revoking token",
		}
	}

	if err := s.refreshTokenRepository.Insert(ctx, &entity.RefreshToken{UserID: user.ID, Token: newRefreshToken}); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed registering token",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "JWT successfully refresh",
		Data: &model.AllToken{
			RefreshToken: newRefreshToken,
			AccessToken:  newAccessToken,
		},
	}
}

func (s *authService) Logout(ctx context.Context, refreshToken string) helpers.BaseResponse {
	cfg := config.AppConfig

	_, err := helpers.ValidateToken(refreshToken, cfg.JwtRefreshPublicSecret)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	tokenEntity, err := s.refreshTokenRepository.FindByToken(ctx, refreshToken)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	user, err := s.userRepository.FindByID(ctx, tokenEntity.UserID)
	if err != nil || user == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	if err := s.refreshTokenRepository.RevokeByToken(ctx, refreshToken); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed revoking token",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("%s user successfully logout", user.Username),
	}
}

func (s *authService) VerifyAccessToken(ctx context.Context, accessToken string) helpers.BaseResponse {
	cfg := config.AppConfig

	_, err := helpers.ValidateToken(accessToken, cfg.JwtAccessPublicSecret)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Access token valid",
	}
}

func (s *authService) VerifyRefreshToken(ctx context.Context, refreshToken string) helpers.BaseResponse {
	cfg := config.AppConfig

	_, err := helpers.ValidateToken(refreshToken, cfg.JwtRefreshPublicSecret)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	tokenEntity, err := s.refreshTokenRepository.FindByToken(ctx, refreshToken)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	user, err := s.userRepository.FindByID(ctx, tokenEntity.UserID)
	if err != nil || user == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Errors:  err,
			Message: "Invalid token",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Refresh token valid",
	}
}
