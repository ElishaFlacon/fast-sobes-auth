package accessToken

import (
	"time"
)

type AccessToken struct {
	Id        int64     `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserId    int64     `gorm:"index;not null"`
	Revoked   bool      `gorm:"default:false;index"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
