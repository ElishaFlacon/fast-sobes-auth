package accessToken

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) toDomain(model *AccessToken) *domain.AccessToken {
	return &domain.AccessToken{
		ID:        model.ID,
		JTI:       model.JTI,
		UserID:    model.UserID,
		Revoked:   model.Revoked,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}
}

func (r *repository) toModel(domain *domain.AccessToken) *AccessToken {
	return &AccessToken{
		ID:        domain.ID,
		JTI:       domain.JTI,
		UserID:    domain.UserID,
		Revoked:   domain.Revoked,
		ExpiresAt: domain.ExpiresAt,
		CreatedAt: domain.CreatedAt,
	}
}
