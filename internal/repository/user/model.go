package user

import (
	"time"
)

type User struct {
	Id               int64  `gorm:"primaryKey"`
	Email            string `gorm:"uniqueIndex;not null"`
	PasswordHash     string `gorm:"not null"`
	PermissionLevel  int32  `gorm:"default:1"`
	Disabled         bool   `gorm:"default:false"`
	TwoFactorEnabled bool   `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (User) TableName() string {
	return "users"
}
