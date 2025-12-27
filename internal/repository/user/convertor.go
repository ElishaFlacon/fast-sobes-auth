package user

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) toDomain(model *User) *domain.User {
	return &domain.User{
		ID:               model.ID,
		Email:            model.Email,
		PasswordHash:     model.PasswordHash,
		PermissionLevel:  model.PermissionLevel,
		Disabled:         model.Disabled,
		TwoFactorEnabled: model.TwoFactorEnabled,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

func (r *repository) toModel(domain *domain.User) *User {
	return &User{
		ID:               domain.ID,
		Email:            domain.Email,
		PasswordHash:     domain.PasswordHash,
		PermissionLevel:  domain.PermissionLevel,
		Disabled:         domain.Disabled,
		TwoFactorEnabled: domain.TwoFactorEnabled,
		CreatedAt:        domain.CreatedAt,
		UpdatedAt:        domain.UpdatedAt,
	}
}
