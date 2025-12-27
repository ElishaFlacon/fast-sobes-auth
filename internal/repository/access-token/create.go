package accessToken

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Create(ctx context.Context, token *domain.AccessToken) error {
	model := r.toModel(token)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	token.ID = model.ID
	return nil
}
