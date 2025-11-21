package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) Login(ctx context.Context, email, password string) (*domain.LoginResult, error) {
	u.log.Infof("Login user: %s", email)

	// Получение пользователя
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Проверка блокировки
	if user.Disabled {
		return nil, fmt.Errorf("user disabled")
	}

	// Проверка пароля
	passwordHash, err := u.passwordRepo.GetPasswordHash(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	result := &domain.LoginResult{}

	// Если включена 2FA
	if user.TwoFactorEnabled {
		tempToken, err := u.tokenRepo.CreateTempToken(
			ctx,
			user.ID,
			"login",
			nil,
			time.Now().Add(5*time.Minute),
		)
		if err != nil {
			return nil, fmt.Errorf("create temp token: %w", err)
		}

		result.RequiresTwoFactor = true
		result.TempToken = tempToken
	} else {
		// Создание токенов доступа
		auth, err := u.createAuthTokens(ctx, user)
		if err != nil {
			return nil, err
		}
		result.Auth = auth
	}

	return result, nil
}

func (u *usecase) VerifyTwoFactor(ctx context.Context, tempToken, code string) (*domain.AuthResult, error) {
	u.log.Infof("Verify 2FA")

	// Проверка временного токена
	userID, _, err := u.tokenRepo.VerifyTempToken(ctx, tempToken, "login")
	if err != nil {
		return nil, fmt.Errorf("verify temp token: %w", err)
	}

	// Получение пользователя
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Проверка кода 2FA
	secret, err := u.twoFactorRepo.GetSecret(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get 2FA secret: %w", err)
	}

	if !u.verify2FACode(secret, code) {
		// Попытка использовать backup код
		if err := u.twoFactorRepo.UseBackupCode(ctx, userID, code); err != nil {
			return nil, fmt.Errorf("invalid 2FA code")
		}
	}

	// Удаление временного токена
	if err := u.tokenRepo.RevokeTempToken(ctx, tempToken); err != nil {
		u.log.Errorf("Failed to revoke temp token: %v", err)
	}

	// Создание токенов доступа
	return u.createAuthTokens(ctx, user)
}
