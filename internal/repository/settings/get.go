package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Get(ctx context.Context) (*domain.Settings, error) {
	var settings Settings

	err := r.db.WithContext(ctx).
		Where("id = ?", defaultSettingsID).
		First(&settings).Error

	return r.toDomain(&settings), err
}
