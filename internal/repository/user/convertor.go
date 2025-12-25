package user

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) toDomain(model *User) *domain.User {
	return &domain.User{
		Id:               model.Id,
		Email:            model.Email,
		PermissionLevel:  model.PermissionLevel,
		Disabled:         model.Disabled,
		TwoFactorEnabled: model.TwoFactorEnabled,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

func (r *repository) toModel(domain *domain.User) *User {
	return &User{
		Id:               domain.Id,
		Email:            domain.Email,
		PermissionLevel:  domain.PermissionLevel,
		Disabled:         domain.Disabled,
		TwoFactorEnabled: domain.TwoFactorEnabled,
		CreatedAt:        domain.CreatedAt,
		UpdatedAt:        domain.UpdatedAt,
	}
}
