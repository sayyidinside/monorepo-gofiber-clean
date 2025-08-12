package database

import (
	"fmt"
	"log"

	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	var dbHost, dbPort, dbUser, dbPassword, dbName string

	cfg := config.AppConfig

	switch cfg.Env {
	case "production":
		dbHost = cfg.ProdDbHost
		dbPort = cfg.ProdDbPort
		dbUser = cfg.ProdDbUsername
		dbPassword = cfg.ProdDbPassword
		dbName = cfg.ProdDbName
	case "development":
		dbHost = cfg.DevDbHost
		dbPort = cfg.DevDbPort
		dbUser = cfg.DevDbUsername
		dbPassword = cfg.DevDbPassword
		dbName = cfg.DevDbName
	default:
		dbHost = cfg.LocalDbHost
		dbPort = cfg.LocalDbPort
		dbUser = cfg.LocalDbUsername
		dbPassword = cfg.LocalDbPassword
		dbName = cfg.LocalDbName
	}

	// Construct database URL
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	return db, nil
}
