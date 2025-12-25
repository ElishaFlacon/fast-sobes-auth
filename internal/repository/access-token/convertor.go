package accessToken

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) toDomain(model *AccessToken) *domain.AccessToken {
	return &domain.AccessToken{
		Id:        model.Id,
		Token:     model.Token,
		UserId:    model.UserId,
		Revoked:   model.Revoked,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}
}

func (r *repository) toModel(domain *domain.AccessToken) *AccessToken {
	return &AccessToken{
		Id:        domain.Id,
		Token:     domain.Token,
		UserId:    domain.UserId,
		Revoked:   domain.Revoked,
		ExpiresAt: domain.ExpiresAt,
		CreatedAt: domain.CreatedAt,
	}
}
