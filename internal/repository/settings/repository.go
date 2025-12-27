package settings

import (
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	"gorm.io/gorm"
)

var _ def.SettingsRepository = (*repository)(nil)

const defaultSettingsID = 0

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Settings{},
	)
}

func (r *repository) getDefaultSettings() Settings {
	return Settings{
		ID:                        defaultSettingsID,
		RequireTwoFactor:          false,
		TokenTTLMinutes:           60,
		MinPasswordLength:         8,
		RequirePasswordComplexity: true,
	}
}
