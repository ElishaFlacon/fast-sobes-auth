package tokens

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) VerifyToken(ctx context.Context, token string) (*domain.TokenInfo, error) {
	u.log.Infof("Verify token")

	tokenInfo, err := u.tokenRepo.VerifyAccessToken(ctx, token)
	if err != nil {
		return &domain.TokenInfo{Valid: false}, nil
	}

	return tokenInfo, nil
}