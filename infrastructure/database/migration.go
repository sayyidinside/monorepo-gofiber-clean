package database

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Module{})
	db.AutoMigrate(&entity.RolePermission{})
	db.AutoMigrate(&entity.Permission{})
	db.AutoMigrate(&entity.Role{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.RefreshToken{})
}
