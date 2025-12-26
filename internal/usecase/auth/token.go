package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/jwtmanager"
)

func (u *usecase) issueAuth(ctx context.Context, user *domain.User) (*domain.AuthResult, error) {
	settings, err := u.loadSettings(ctx)
	if err != nil {
		return nil, err
	}

	now := u.now()
	expiresAt := now.Add(time.Duration(settings.TokenTTLMinutes) * time.Minute)

	token, err := u.jwt.Sign(jwtmanager.Claims{
		UserID:          user.Id,
		Email:           user.Email,
		PermissionLevel: user.PermissionLevel,
		ExpiresAt:       expiresAt.Unix(),
		IssuedAt:        now.Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	if err := u.tokens.Create(ctx, &domain.AccessToken{
		Token:     token,
		UserId:    user.Id,
		Revoked:   false,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}); err != nil {
		return nil, fmt.Errorf("store token: %w", err)
	}

	u.log.Infof("Issued token for user_id=%d exp=%s", user.Id, expiresAt.Format(time.RFC3339))

	return &domain.AuthResult{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        def.SanitizeUser(user),
	}, nil
}
