package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) UpdatePermissionLevel(ctx context.Context, id string, level int32) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("permission_level", level).Error
}

func (r *repository) SetDisabled(ctx context.Context, id string, disabled bool) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("disabled", disabled).Error
}

func (r *repository) UpdateTwoFactorSecret(ctx context.Context, id, secret string) error {
	// Реализовано в TwoFactorRepository
	return nil
}

func (r *repository) UpdateTwoFactorEnabled(ctx context.Context, id string, enabled bool) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("two_factor_enabled", enabled).Error
}

func (r *repository) UpdateEmail(ctx context.Context, id, newEmail string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("email", newEmail).Error
}

func (r *repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *repository) toDomain(model *User) *domain.User {
	return &domain.User{
		ID:               model.ID,
		Email:            model.Email,
		PermissionLevel:  model.PermissionLevel,
		Disabled:         model.Disabled,
		TwoFactorEnabled: model.TwoFactorEnabled,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}
