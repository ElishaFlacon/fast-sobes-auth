package tokens

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResult, error) {
	u.log.Infof("Refresh token")

	// Проверка refresh token
	userId, err := u.tokenRepo.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Получение пользователя
	user, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Проверка блокировки
	if user.Disabled {
		return nil, fmt.Errorf("user disabled")
	}

	// Отзыв старого refresh token
	if err := u.tokenRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		u.log.Errorf("Failed to revoke refresh token: %v", err)
	}

	// Создание новых токенов
	return u.createAuthTokens(ctx, user)
}
