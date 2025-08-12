package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID   uint      `json:"id" gorm:"primaryKey"`
	UUID uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	Name string    `json:"name" gorm:"not null"`

	// Relationship
	Permissions []Permission `json:"permissions" gorm:"foreignKey:ModuleID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Module) TableName() string {
	return "modules"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (m *Module) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return
}
