package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email          string `gorm:"type:varchar(20);unique;index" json:"email"`
	PhoneNumber    string `gorm:"type:varchar(20);unique" json:"phone_number"`
	HashedPassword string
	LastLogin      time.Time
	FirstName      string
	LastName       string
	IsVerified     bool `gorm:"default:false"`
	IsSuperUser    bool `gorm:"default:false"`
}

type BlacklistedToken struct {
	Token     string `gorm:"primary_key"`
	ExpiresAt time.Time
}

type EmailVerificationToken struct {
	UUID      uuid.UUID `gorm:"type:uuid;primary_key" json:"uuid"`
	ExpiresAt time.Time
	UserID    uint
	User      User
}
type PasswordResetToken struct {
	UUID      uuid.UUID `gorm:"type:uuid;primary_key" json:"uuid"`
	ExpiresAt time.Time
	UserID    uint
	User      User
}
