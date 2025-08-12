package database

import (
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Module{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.RolePermission{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Permission{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Role{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.RefreshToken{}); err != nil {
		return err
	}

	return nil
}
