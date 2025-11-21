package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) Register(ctx context.Context, email, password string) (*domain.RegisterResult, error) {
	u.log.Infof("Register user: %s", email)

	// Проверка существования пользователя
	exists, err := u.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check user exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Получение настроек
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	// Валидация пароля
	if err := u.validatePassword(password, settings); err != nil {
		return nil, err
	}

	// Создание пользователя
	user := &domain.User{
		ID:               uuid.New().String(),
		Email:            email,
		PermissionLevel:  0,
		Disabled:         false,
		TwoFactorEnabled: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	if err := u.passwordRepo.SetPassword(ctx, user.ID, string(hash)); err != nil {
		return nil, fmt.Errorf("save password: %w", err)
	}

	result := &domain.RegisterResult{
		User:              user,
		RequiresTwoFactor: settings.RequireTwoFactor,
	}

	// Если требуется 2FA
	if settings.RequireTwoFactor {
		secret, qrCode, err := u.generate2FASecret(user.Email)
		if err != nil {
			return nil, fmt.Errorf("generate 2FA secret: %w", err)
		}

		// Сохранение временного секрета
		if err := u.twoFactorRepo.SaveSecret(ctx, user.ID, secret); err != nil {
			return nil, fmt.Errorf("save 2FA secret: %w", err)
		}

		// Создание временного токена
		tempToken, err := u.tokenRepo.CreateTempToken(
			ctx,
			user.ID,
			"registration",
			nil,
			time.Now().Add(15*time.Minute),
		)
		if err != nil {
			return nil, fmt.Errorf("create temp token: %w", err)
		}

		result.TwoFactorSecret = secret
		result.QRCodeURL = qrCode
		result.TempToken = tempToken
	}

	return result, nil
}

func (u *usecase) VerifyTwoFactorSetup(ctx context.Context, tempToken, code string) (*domain.AuthResult, error) {
	u.log.Infof("Verify 2FA setup")

	// Проверка временного токена
	userID, _, err := u.tokenRepo.VerifyTempToken(ctx, tempToken, "registration")
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
		return nil, fmt.Errorf("invalid 2FA code")
	}

	// Активация 2FA
	user.TwoFactorEnabled = true
	if err := u.userRepo.UpdateTwoFactorEnabled(ctx, userID, true); err != nil {
		return nil, fmt.Errorf("enable 2FA: %w", err)
	}

	// Удаление временного токена
	if err := u.tokenRepo.RevokeTempToken(ctx, tempToken); err != nil {
		u.log.Errorf("Failed to revoke temp token: %v", err)
	}

	// Создание токенов доступа
	return u.createAuthTokens(ctx, user)
}

func (u *usecase) validatePassword(password string, settings *domain.Settings) error {
	if len(password) < int(settings.MinPasswordLength) {
		return fmt.Errorf("password too short")
	}

	if settings.RequirePasswordComplexity {
		// Проверка сложности пароля
		// TODO: реализовать проверку на наличие букв, цифр, спецсимволов
	}

	return nil
}

func (u *usecase) generate2FASecret(email string) (secret, qrCode string, err error) {
	// TODO: реализовать генерацию секрета и QR-кода с использованием библиотеки TOTP
	return "", "", nil
}

func (u *usecase) verify2FACode(secret, code string) bool {
	// TODO: реализовать проверку TOTP кода
	return false
}

func (u *usecase) createAuthTokens(ctx context.Context, user *domain.User) (*domain.AuthResult, error) {
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	// Создание access token
	accessToken, err := u.tokenRepo.CreateAccessToken(
		ctx,
		user.ID,
		time.Now().Add(time.Duration(settings.TokenTTLMinutes)*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	// Создание refresh token
	refreshToken, err := u.tokenRepo.CreateRefreshToken(
		ctx,
		user.ID,
		time.Now().Add(time.Duration(settings.RefreshTokenTTLDays)*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	return &domain.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(settings.TokenTTLMinutes * 60),
		User:         user,
	}, nil
}
