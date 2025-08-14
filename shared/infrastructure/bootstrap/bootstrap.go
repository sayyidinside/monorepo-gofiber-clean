package bootstrap

import (
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/database"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/redis"
	"gorm.io/gorm"
)

type Deps struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.RedisClient
}

func NewDeps(env string) (*Deps, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, err
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	r := redis.Connect(cfg)

	return &Deps{Config: cfg, DB: db, Redis: r}, nil
}
