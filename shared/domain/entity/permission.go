package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UUID     uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	Name     string    `json:"name" gorm:"not null"`
	ModuleID uint      `json:"module_id" gorm:"not null"`

	// Relationship
	Module Module `json:"module" gorm:"foreignKey:ModuleID"`
	Roles  []Role `json:"roles" gorm:"many2many:role_permissions;"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UUID == uuid.Nil {
		p.UUID = uuid.New()
	}
	return
}
