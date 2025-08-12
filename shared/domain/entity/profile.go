package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID         uint      `json:"id"`
	UUID       uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	UserID     uint      `json:"user_id" gorm:"index;not null"`
	Name       string    `json:"name" gorm:"not null"`
	BirthPlace string    `json:"birth_place" gorm:"not null"`
	BirthDate  time.Time `json:"birth_date" gorm:"type:date;not null"`
	Address    string    `json:"address"`

	// Relationship
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Profile) TableName() string {
	return "profiles"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field.
func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UUID == uuid.Nil {
		p.UUID = uuid.New()
	}

	return
}

func (p *Profile) GetID() uint {
	if p == nil {
		return 0
	}

	return p.ID
}
