package repository

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Password{},
		&models.AccessToken{},
		&models.RefreshToken{},
		&models.TempToken{},
		&models.TwoFactorSecret{},
		&models.BackupCode{},
		&models.Settings{},
		&models.EmailChangeRequest{},
	)
}
