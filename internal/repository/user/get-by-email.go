package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}
