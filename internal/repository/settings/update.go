package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Update(ctx context.Context, settings *domain.Settings) error {
	model := r.toModel(settings)

	err := r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", defaultSettingsId).
		Updates(model).Error

	return err
}

func (r *repository) UpdateById(ctx context.Context, id int64, settings *domain.Settings) error {
	model := r.toModel(settings)

	err := r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", id).
		Updates(model).Error

	return err
}
