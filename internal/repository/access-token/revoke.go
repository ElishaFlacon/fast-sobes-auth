package accessToken

import "context"

func (r *repository) Revoke(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&AccessToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *repository) RevokeAllByUser(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).
		Model(&AccessToken{}).
		Where("user_id = ?", userId).
		Update("revoked", true).Error
}
