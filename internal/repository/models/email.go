package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailChangeRequest struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	NewEmail  string    `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}

func (e *EmailChangeRequest) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

func (EmailChangeRequest) TableName() string {
	return "email_change_requests"
}
