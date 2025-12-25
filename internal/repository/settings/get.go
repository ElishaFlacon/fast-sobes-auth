package settings

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) Get(ctx context.Context) (*domain.Settings, error) {
	var settings Settings

	err := r.db.WithContext(ctx).
		Where("id = ?", defaultSettingsId).
		First(&settings).Error

	return r.toDomain(&settings), err
}

func (r *repository) GetById(ctx context.Context, id int64) (*domain.Settings, error) {
	var settings Settings

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&settings).Error

	return r.toDomain(&settings), err
}
