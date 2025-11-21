package models

import (
	"time"
)

type TwoFactorSecret struct {
	UserID    string `gorm:"type:uuid;primaryKey"`
	Secret    string `gorm:"not null"`
	UpdatedAt time.Time
}

func (TwoFactorSecret) TableName() string {
	return "two_factor_secrets"
}

type BackupCode struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"type:uuid;index;not null"`
	Code      string `gorm:"uniqueIndex;not null"`
	Used      bool   `gorm:"default:false"`
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (BackupCode) TableName() string {
	return "backup_codes"
}
