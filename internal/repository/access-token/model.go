package accessToken

import (
	"time"
)

type AccessToken struct {
	ID        int64     `gorm:"primaryKey"`
	JTI       string    `gorm:"uniqueIndex;not null"`
	UserID    int64     `gorm:"index;not null"`
	Revoked   bool      `gorm:"default:false;index"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
