package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	var model User

	if err := r.db.WithContext(ctx).
		Where(&User{Id: id}).
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}
