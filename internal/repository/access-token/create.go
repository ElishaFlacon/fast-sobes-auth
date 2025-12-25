package accessToken

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Create(ctx context.Context, token *domain.AccessToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}
