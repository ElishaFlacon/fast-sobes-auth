package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	model := &User{
		ID:               user.ID,
		Email:            user.Email,
		PermissionLevel:  user.PermissionLevel,
		Disabled:         user.Disabled,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}

	return r.db.WithContext(ctx).Create(model).Error
}
