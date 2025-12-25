package settings

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) toDomain(model *Settings) *domain.Settings {
	return &domain.Settings{
		Id:                        model.Id,
		RequireTwoFactor:          model.RequireTwoFactor,
		TokenTTLMinutes:           model.TokenTTLMinutes,
		MinPasswordLength:         model.MinPasswordLength,
		RequirePasswordComplexity: model.RequirePasswordComplexity,
		UpdatedAt:                 model.UpdatedAt,
	}
}

func (r *repository) toModel(domain *domain.Settings) *Settings {
	return &Settings{
		Id:                        domain.Id,
		RequireTwoFactor:          domain.RequireTwoFactor,
		TokenTTLMinutes:           domain.TokenTTLMinutes,
		MinPasswordLength:         domain.MinPasswordLength,
		RequirePasswordComplexity: domain.RequirePasswordComplexity,
		UpdatedAt:                 domain.UpdatedAt,
	}
}
