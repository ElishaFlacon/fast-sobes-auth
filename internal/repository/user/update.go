package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Update(ctx context.Context, user *domain.User) error {
	model := r.toModel(user)

	return r.db.WithContext(ctx).
		Model(&User{}).
		Where(&User{Id: user.Id}).
		Updates(model).Error
}
