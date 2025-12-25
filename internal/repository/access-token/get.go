package accessToken

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) GetByToken(ctx context.Context, token string) (*domain.AccessToken, error) {
	var accessToken AccessToken

	if err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&accessToken).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&accessToken), nil
}
