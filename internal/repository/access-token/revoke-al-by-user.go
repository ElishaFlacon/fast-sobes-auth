package accessToken

import "context"

func (r *repository) RevokeAllByUser(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&AccessToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}
