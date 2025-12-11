package settings

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) UpdateSettings(
	ctx context.Context,
	requireTwoFactor *bool,
	tokenTTLMinutes *int32,
	refreshTokenTTLDays *int32,
	minPasswordLength *int32,
	requirePasswordComplexity *bool,
) (*domain.Settings, error) {
	u.log.Infof("Update settings")

	// Получение текущих настроек
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	// Обновление полей
	if requireTwoFactor != nil {
		settings.RequireTwoFactor = *requireTwoFactor
	}
	if tokenTTLMinutes != nil {
		if *tokenTTLMinutes < 1 || *tokenTTLMinutes > 1440 {
			return nil, fmt.Errorf("invalid token TTL")
		}
		settings.TokenTTLMinutes = *tokenTTLMinutes
	}
	if refreshTokenTTLDays != nil {
		if *refreshTokenTTLDays < 1 || *refreshTokenTTLDays > 365 {
			return nil, fmt.Errorf("invalid refresh token TTL")
		}
		settings.RefreshTokenTTLDays = *refreshTokenTTLDays
	}
	if minPasswordLength != nil {
		if *minPasswordLength < 6 || *minPasswordLength > 128 {
			return nil, fmt.Errorf("invalid min password length")
		}
		settings.MinPasswordLength = *minPasswordLength
	}
	if requirePasswordComplexity != nil {
		settings.RequirePasswordComplexity = *requirePasswordComplexity
	}

	// Сохранение настроек
	if err := u.settingsRepo.Update(ctx, settings); err != nil {
		return nil, fmt.Errorf("update settings: %w", err)
	}

	return settings, nil
}
