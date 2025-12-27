package accessToken

import "context"

func (r *repository) Revoke(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).
		Model(&AccessToken{}).
		Where(&AccessToken{JTI: jti}).
		Update("revoked", true).Error
}
