package accessToken

import (
	"context"
	"time"
)

func (r *repository) DeleteExpired(ctx context.Context, now time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&AccessToken{}).Error
}
