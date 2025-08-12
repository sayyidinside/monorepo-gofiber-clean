package entity

import (
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UUID        uuid.UUID    `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	RoleID      uint         `json:"role_id"`
	Username    string       `json:"username" gorm:"index"`
	Email       string       `json:"email" gorm:"index"`
	Password    string       `json:"password"`
	ValidatedAt sql.NullTime `json:"validated_at" gorm:"index"`
	Role        Role         `json:"role" gorm:"foreignKey:RoleID"`
	gorm.Model
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	return
}

func (u *User) GetID() uint {
	if u == nil {
		return 0
	}

	return u.ID
}
