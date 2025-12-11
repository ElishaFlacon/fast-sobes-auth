package auth

import (
	"context"
)

func (u *usecase) Logout(ctx context.Context, token, refreshToken string) error {
	u.log.Infof("Logout")

	// Отзыв access token
	if token != "" {
		if err := u.tokenRepo.RevokeAccessToken(ctx, token); err != nil {
			u.log.Errorf("Failed to revoke access token: %v", err)
		}
	}

	// Отзыв refresh token
	if refreshToken != "" {
		if err := u.tokenRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
			u.log.Errorf("Failed to revoke refresh token: %v", err)
		}
	}

	return nil
}
