package settings

import (
	"context"
)

func (r *repository) Reset(ctx context.Context) error {
	settings := r.getDefaultSettings()

	return r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", defaultSettingsId).
		Updates(&settings).Error
}

func (r *repository) ResetById(ctx context.Context, id int64) error {
	settings := r.getDefaultSettings()

	return r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", id).
		Updates(&settings).Error
}
