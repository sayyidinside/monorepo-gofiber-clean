package service

import (
	"context"
	"encoding/json"
	"log"

	sharedRepository "github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
)

type EmailService interface {
	SendEmail(ctx context.Context, email model.Email) error
}

type emailService struct {
	userRepository sharedRepository.UserRepository
	cacheRedis     *redis.CacheClient
}

func NewEmailService(cacheRedis *redis.CacheClient, userRepository sharedRepository.UserRepository) EmailService {
	return &emailService{
		userRepository: userRepository,
		cacheRedis:     cacheRedis,
	}
}

func (s *emailService) SendEmail(ctx context.Context, email model.Email) error {
	// TODO: Implement SMTP Email
	jsonData, _ := json.MarshalIndent(email, "", " ")
	log.Println(string(jsonData))
	return nil
}
