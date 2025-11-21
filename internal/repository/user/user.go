package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
)

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	model := &models.User{
		ID:               user.ID,
		Email:            user.Email,
		PermissionLevel:  user.PermissionLevel,
		Disabled:         user.Disabled,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}

	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model models.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model models.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *repository) Update(ctx context.Context, user *domain.User) error {
	model := &models.User{
		ID:               user.ID,
		Email:            user.Email,
		PermissionLevel:  user.PermissionLevel,
		Disabled:         user.Disabled,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}

	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(model).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *repository) List(
	ctx context.Context,
	offset,
	limit int32,
	minPermissionLevel *int32,
	includeDisabled bool,
) ([]*domain.User, int32, error) {
	var modelUsers []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if minPermissionLevel != nil {
		query = query.Where("permission_level >= ?", *minPermissionLevel)
	}

	if !includeDisabled {
		query = query.Where("disabled = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&modelUsers).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(modelUsers))
	for i, model := range modelUsers {
		users[i] = r.toDomain(&model)
	}

	return users, int32(total), nil
}

func (r *repository) UpdatePermissionLevel(ctx context.Context, id string, level int32) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("permission_level", level).Error
}

func (r *repository) SetDisabled(ctx context.Context, id string, disabled bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("disabled", disabled).Error
}

func (r *repository) UpdateTwoFactorSecret(ctx context.Context, id, secret string) error {
	// Реализовано в TwoFactorRepository
	return nil
}

func (r *repository) UpdateTwoFactorEnabled(ctx context.Context, id string, enabled bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("two_factor_enabled", enabled).Error
}

func (r *repository) UpdateEmail(ctx context.Context, id, newEmail string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("email", newEmail).Error
}

func (r *repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *repository) toDomain(model *models.User) *domain.User {
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
