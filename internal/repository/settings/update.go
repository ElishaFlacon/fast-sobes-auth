package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Update(ctx context.Context, settings *domain.Settings) error {
	model := r.toModel(settings)

	err := r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", defaultSettingsID).
		Updates(model).Error

	return err
}
