package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
	"gorm.io/gorm"
)

func (r *repository) Get(ctx context.Context) (*domain.Settings, error) {
	var settings models.Settings

	err := r.db.WithContext(ctx).Where("id = ?", defaultSettingsID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// Создание настроек по умолчанию
		settings = r.getDefaultSettings()
		if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&settings), nil
}

func (r *repository) Update(ctx context.Context, settings *domain.Settings) error {
	model := &models.Settings{
		ID:                        defaultSettingsID,
		RequireTwoFactor:          settings.RequireTwoFactor,
		TokenTTLMinutes:           settings.TokenTTLMinutes,
		RefreshTokenTTLDays:       settings.RefreshTokenTTLDays,
		MinPasswordLength:         settings.MinPasswordLength,
		RequirePasswordComplexity: settings.RequirePasswordComplexity,
	}

	return r.db.WithContext(ctx).Model(&models.Settings{}).
		Where("id = ?", defaultSettingsID).
		Updates(model).Error
}

func (r *repository) Reset(ctx context.Context) error {
	settings := r.getDefaultSettings()

	return r.db.WithContext(ctx).Model(&models.Settings{}).
		Where("id = ?", defaultSettingsID).
		Updates(&settings).Error
}

func (r *repository) getDefaultSettings() models.Settings {
	return models.Settings{
		ID:                        defaultSettingsID,
		RequireTwoFactor:          false,
		TokenTTLMinutes:           60,
		RefreshTokenTTLDays:       30,
		MinPasswordLength:         8,
		RequirePasswordComplexity: true,
	}
}

func (r *repository) toDomain(model *models.Settings) *domain.Settings {
	return &domain.Settings{
		ID:                        model.ID,
		RequireTwoFactor:          model.RequireTwoFactor,
		TokenTTLMinutes:           model.TokenTTLMinutes,
		RefreshTokenTTLDays:       model.RefreshTokenTTLDays,
		MinPasswordLength:         model.MinPasswordLength,
		RequirePasswordComplexity: model.RequirePasswordComplexity,
		UpdatedAt:                 model.UpdatedAt,
	}
}
