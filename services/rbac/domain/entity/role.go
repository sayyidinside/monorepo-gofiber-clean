package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	UUID    uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	Name    string    `json:"name" gorm:"size:50;uniqueIndex;not null"`
	IsAdmin bool      `json:"is_admin" gorm:"default:false"`

	// Relationship
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users" gorm:"foreignKey:RoleID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.UUID == uuid.Nil {
		r.UUID = uuid.New()
	}
	return
}
