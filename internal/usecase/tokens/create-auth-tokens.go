package tokens

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) createAuthTokens(ctx context.Context, user *domain.User) (*domain.AuthResult, error) {
	// Генерируем токены в usecase
	accessToken, err := u.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := u.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// Вычисляем время истечения
	accessExpiresAt := time.Now().Add(u.cfg.AccessTokenTTL)
	refreshExpiresAt := time.Now().Add(u.cfg.RefreshTokenTTL)

	// Сохраняем через репозиторий
	if err := u.tokenRepo.SaveAccessToken(ctx, accessToken, user.ID, accessExpiresAt); err != nil {
		return nil, fmt.Errorf("save access token: %w", err)
	}

	if err := u.tokenRepo.SaveRefreshToken(ctx, refreshToken, user.ID, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &domain.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.cfg.AccessTokenTTL.Seconds()),
		User:         user,
	}, nil
}

func (u *usecase) generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
