package bootstrap

import (
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/database"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/rabbitmq"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"gorm.io/gorm"
)

type Deps struct {
	Config   *config.Config
	DB       *gorm.DB
	Redis    *redis.RedisClient
	RabbitMQ *rabbitmq.RabbitMQClient
}

func NewDeps(env string, is_rabbitmq bool) (*Deps, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, err
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	r := redis.Connect(cfg)

	mq := &rabbitmq.RabbitMQClient{}
	if is_rabbitmq {
		mq, err = rabbitmq.Connect(config.AppConfig)
		if err != nil {
			return nil, err
		}
	} else {
		mq = nil
	}

	return &Deps{Config: cfg, DB: db, Redis: r, RabbitMQ: mq}, nil
}
