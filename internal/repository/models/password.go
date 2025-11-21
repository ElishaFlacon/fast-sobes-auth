package models

import (
	"time"
)

type Password struct {
	UserID       string `gorm:"type:uuid;primaryKey"`
	PasswordHash string `gorm:"not null"`
	UpdatedAt    time.Time
}

func (Password) TableName() string {
	return "passwords"
}
