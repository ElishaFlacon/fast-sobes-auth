package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"gorm.io/gorm"
)

func (u *usecase) loadSettings(ctx context.Context) (*domain.Settings, error) {
	settings, err := u.settings.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.log.Infof("Settings not found, using defaults")
			return def.DefaultSettings(), nil
		}
		return nil, fmt.Errorf("load settings: %w", err)
	}

	return settings, nil
}
