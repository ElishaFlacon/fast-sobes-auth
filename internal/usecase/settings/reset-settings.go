package settings

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) ResetSettings(ctx context.Context) (*domain.Settings, error) {
	u.log.Infof("Reset settings")

	if err := u.settingsRepo.Reset(ctx); err != nil {
		return nil, fmt.Errorf("reset settings: %w", err)
	}

	return u.settingsRepo.Get(ctx)
}
