package settings

import (
	"context"
)

func (r *repository) Reset(ctx context.Context) error {
	settings := r.getDefaultSettings()

	return r.db.WithContext(ctx).Model(&Settings{}).
		Where("id = ?", defaultSettingsID).
		Updates(&settings).Error
}
