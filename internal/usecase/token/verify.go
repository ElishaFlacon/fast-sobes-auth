package tokens

import (
	"context"
	"errors"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"gorm.io/gorm"
)

func (u *usecase) Verify(ctx context.Context, token string) (*domain.TokenInfo, error) {
	u.log.Infof("Verify token")

	claims, err := u.jwt.Verify(token)
	if err != nil {
		return nil, fmt.Errorf("verify jwt: %w", err)
	}

	stored, err := u.tokens.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("load token: %w", err)
	}

	now := u.now()

	if stored.Revoked {
		return &domain.TokenInfo{
			Token:     token,
			UserID:    stored.UserID,
			ExpiresAt: stored.ExpiresAt,
			Revoked:   true,
			Valid:     false,
		}, fmt.Errorf("token revoked")
	}

	if now.After(stored.ExpiresAt) || claims.ExpiresAt < now.Unix() {
		return &domain.TokenInfo{
			Token:     token,
			UserID:    stored.UserID,
			ExpiresAt: stored.ExpiresAt,
			Revoked:   stored.Revoked,
			Valid:     false,
		}, fmt.Errorf("token expired")
	}

	if stored.UserID != claims.UserID {
		return &domain.TokenInfo{
			Token:     token,
			UserID:    stored.UserID,
			ExpiresAt: stored.ExpiresAt,
			Revoked:   stored.Revoked,
			Valid:     false,
		}, fmt.Errorf("token user mismatch")
	}

	user, err := u.users.GetByID(ctx, stored.UserID)
	if err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}

	return &domain.TokenInfo{
		Token:     token,
		UserID:    stored.UserID,
		User:      def.SanitizeUser(user),
		ExpiresAt: stored.ExpiresAt,
		Revoked:   stored.Revoked,
		Valid:     true,
	}, nil
}
