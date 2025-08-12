package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/services/profile/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	sharedRepository "github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type ProfileService interface {
	GetByUserID(ctx context.Context, id uint) helpers.BaseResponse
	GetByUserUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	UpdateByUserID(ctx context.Context, input *model.ProfileInput, id uint) helpers.BaseResponse
	UpdateByUserUUID(ctx context.Context, input *model.ProfileInputByAdmin, uuid uuid.UUID) helpers.BaseResponse
}

type profileService struct {
	profileRepository repository.ProfileRepository
	userRepository    sharedRepository.UserRepository
	cacheRedis        *redis.CacheClient
}

func NewProfileService(
	profileRepository repository.ProfileRepository,
	userRepository sharedRepository.UserRepository,
	cacheRedis *redis.CacheClient,
) ProfileService {
	return &profileService{
		profileRepository: profileRepository,
		userRepository:    userRepository,
		cacheRedis:        cacheRedis,
	}
}

func (s *profileService) GetByUserID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	profile, err := s.profileRepository.FindByUserID(ctx, id)
	if profile == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
		})
	}

	profileModel := model.ProfileToDetailModel(profile)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Profile data found",
		Data:    profileModel,
	}
}

func (s *profileService) GetByUserUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
		})
	}

	profile, err := s.profileRepository.FindByUserID(ctx, user.ID)
	if profile == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
		})
	}

	profileModel := model.ProfileToDetailModel(profile)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Profile data found",
		Data:    profileModel,
	}
}

func (s *profileService) UpdateByUserID(ctx context.Context, input *model.ProfileInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check existence
	profile, err := s.profileRepository.FindByUserID(ctx, id)
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
			Errors:  err,
		})
	}

	profileEntity, err := input.ToEntity()
	if profileEntity == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
			Errors:  err,
		})
	}

	profileEntity.UserID = id
	if profile != nil {
		profileEntity.ID = profile.ID
	}

	if err := s.validateEntityInput(ctx, profileEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if profile == nil {
		if err := s.profileRepository.Insert(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	} else {
		if err := s.profileRepository.Update(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Profile succeessfully updated",
	})
}

func (s *profileService) UpdateByUserUUID(ctx context.Context, input *model.ProfileInputByAdmin, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check existence
	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
			Errors:  err,
		})
	}

	profile, err := s.profileRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Profile not found",
			Errors:  err,
		})
	}

	profileEntity, err := input.ToEntity()
	if profileEntity == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
			Errors:  err,
		})
	}

	profileEntity.UserID = user.ID
	if profile != nil {
		profileEntity.ID = profile.ID
	}

	if err := s.validateEntityInput(ctx, profileEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if profile == nil {
		if err := s.profileRepository.Insert(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	} else {
		if err := s.profileRepository.Update(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Profile succeessfully updated",
	})
}

func (s *profileService) validateEntityInput(ctx context.Context, profile *entity.Profile) any {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errs := []helpers.ValidationError{}

	// Check user id valid
	if profile.UserID != 0 {
		if user, err := s.userRepository.FindByID(ctx, profile.UserID); user == nil || err != nil {
			errs = append(errs, helpers.ValidationError{
				Field: "user_id",
				Tag:   "invalid",
			})
		}
	}

	if len(errs) != 0 {
		logData.Message = "Validation error"
		logData.Err = errs
		return errs
	}

	return nil
}
