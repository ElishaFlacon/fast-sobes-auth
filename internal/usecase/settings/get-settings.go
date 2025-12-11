package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) GetSettings(ctx context.Context) (*domain.Settings, error) {
	u.log.Infof("Get settings")

	return u.settingsRepo.Get(ctx)
}
