package settings

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) GetSettings(ctx context.Context) (*domain.Settings, error) {
	u.log.Infof("Get settings")

	return u.settingsRepo.Get(ctx)
}

func (u *usecase) UpdateSettings(
	ctx context.Context,
	req *domain.UpdateSettingsRequest,
) (*domain.Settings, error) {
	u.log.Infof("Update settings")

	// Получение текущих настроек
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	// Обновление полей
	if req.RequireTwoFactor != nil {
		settings.RequireTwoFactor = *req.RequireTwoFactor
	}
	if req.TokenTTLMinutes != nil {
		if *req.TokenTTLMinutes < 1 || *req.TokenTTLMinutes > 1440 {
			return nil, fmt.Errorf("invalid token TTL")
		}
		settings.TokenTTLMinutes = *req.TokenTTLMinutes
	}
	if req.RefreshTokenTTLDays != nil {
		if *req.RefreshTokenTTLDays < 1 || *req.RefreshTokenTTLDays > 365 {
			return nil, fmt.Errorf("invalid refresh token TTL")
		}
		settings.RefreshTokenTTLDays = *req.RefreshTokenTTLDays
	}
	if req.MinPasswordLength != nil {
		if *req.MinPasswordLength < 6 || *req.MinPasswordLength > 128 {
			return nil, fmt.Errorf("invalid min password length")
		}
		settings.MinPasswordLength = *req.MinPasswordLength
	}
	if req.RequirePasswordComplexity != nil {
		settings.RequirePasswordComplexity = *req.RequirePasswordComplexity
	}

	// Сохранение настроек
	if err := u.settingsRepo.Update(ctx, settings); err != nil {
		return nil, fmt.Errorf("update settings: %w", err)
	}

	return settings, nil
}

func (u *usecase) ResetSettings(ctx context.Context) (*domain.Settings, error) {
	u.log.Infof("Reset settings")

	if err := u.settingsRepo.Reset(ctx); err != nil {
		return nil, fmt.Errorf("reset settings: %w", err)
	}

	return u.settingsRepo.Get(ctx)
}
