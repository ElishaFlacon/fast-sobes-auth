package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID               string    `gorm:"type:uuid;primaryKey"`
    Email            string    `gorm:"uniqueIndex;not null"`
    PermissionLevel  int32     `gorm:"default:0"`
    Disabled         bool      `gorm:"default:false"`
    TwoFactorEnabled bool      `gorm:"default:false"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}

func (User) TableName() string {
    return "users"
}
