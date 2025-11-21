package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	ExpiresAt time.Time `gorm:"index"`
	Revoked   bool      `gorm:"default:false;index"`
	CreatedAt time.Time
}

func (t *AccessToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (AccessToken) TableName() string {
	return "access_tokens"
}

type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	ExpiresAt time.Time `gorm:"index"`
	Revoked   bool      `gorm:"default:false;index"`
	CreatedAt time.Time
}

func (t *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type TempToken struct {
	ID        string                 `gorm:"type:uuid;primaryKey"`
	Token     string                 `gorm:"uniqueIndex;not null"`
	UserID    string                 `gorm:"type:uuid;index;not null"`
	TokenType string                 `gorm:"index;not null"`
	Data      map[string]interface{} `gorm:"type:jsonb"`
	ExpiresAt time.Time              `gorm:"index"`
	Revoked   bool                   `gorm:"default:false;index"`
	CreatedAt time.Time
}

func (t *TempToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (TempToken) TableName() string {
	return "temp_tokens"
}
