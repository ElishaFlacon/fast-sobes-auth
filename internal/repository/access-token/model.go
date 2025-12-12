package accessToken

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	Revoked   bool      `gorm:"default:false;index"`
	ExpiresAt time.Time `gorm:"index"`
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
