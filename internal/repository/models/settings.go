package models

import (
	"time"
)

type Settings struct {
	ID                        string `gorm:"type:uuid;primaryKey"`
	RequireTwoFactor          bool   `gorm:"default:false"`
	TokenTTLMinutes           int32  `gorm:"default:60"`
	RefreshTokenTTLDays       int32  `gorm:"default:30"`
	MinPasswordLength         int32  `gorm:"default:8"`
	RequirePasswordComplexity bool   `gorm:"default:true"`
	UpdatedAt                 time.Time
}

func (Settings) TableName() string {
	return "settings"
}
