package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/repository"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/rabbitmq"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/pkg/helpers"
)

type BroadcastService interface {
	SendEmail(ctx context.Context, user_id uint, subject string, messages string) error
}

type broadcastService struct {
	userRepository repository.UserRepository
	rabbitMQClient *rabbitmq.RabbitMQClient
	cacheRedis     *redis.CacheClient
}

func NewBroadcastService(userRepository repository.UserRepository, rabbitMQClient *rabbitmq.RabbitMQClient, cacheRedis *redis.CacheClient) BroadcastService {
	return &broadcastService{
		userRepository: userRepository,
		rabbitMQClient: rabbitMQClient,
		cacheRedis:     cacheRedis,
	}
}

func (s *broadcastService) SendEmail(ctx context.Context, user_id uint, subject string, messages string) error {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Check user existence
	user := &entity.User{}
	userCacheKey := fmt.Sprintf("cache:user-detail:user-id:%d", user_id)

	if err := s.cacheRedis.GetObject(ctx, userCacheKey, user); err != nil || user.GetID() == 0 {
		foundUser, err := s.userRepository.FindByID(ctx, user_id)
		if foundUser == nil || err != nil {
			return errors.New("user not found")
		}
	}

	// Construct rabbitmq body
	email := model.Email{
		User_id: user_id,
		Subject: subject,
		Content: messages,
	}
	body, err := json.Marshal(email)
	if err != nil {
		log.Error(err)
		return errors.New("error consturcting email structure")
	}

	// Publish email via rabbitmq
	if err := s.rabbitMQClient.Channel.PublishWithContext(
		ctx,
		"",      // Exchange
		"email", // Key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	); err != nil {
		log.Error(err)
		return errors.New("error publishing email to rabbitmq")
	}

	return nil
}
