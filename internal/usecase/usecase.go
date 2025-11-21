package usecase

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

type HelloUsecase interface {
	Hello() string
}

type AuthUsecase interface {
	// Регистрация
	Register(ctx context.Context, email, password string) (*domain.RegisterResult, error)
	VerifyTwoFactorSetup(ctx context.Context, tempToken, code string) (*domain.AuthResult, error)

	// Авторизация
	Login(ctx context.Context, email, password string) (*domain.LoginResult, error)
	VerifyTwoFactor(ctx context.Context, tempToken, code string) (*domain.AuthResult, error)

	// Токены
	VerifyToken(ctx context.Context, token string) (*domain.TokenInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResult, error)
	Logout(ctx context.Context, token, refreshToken string) error

	// Управление профилем
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ChangeEmail(ctx context.Context, userID, newEmail, password string) error
	VerifyEmailChange(ctx context.Context, token string) error

	// Управление 2FA
	EnableTwoFactor(ctx context.Context, userID, password string) (*domain.TwoFactorSetup, error)
	DisableTwoFactor(ctx context.Context, userID, password, code string) error

	// Управление пользователями (админ)
	UpdatePermissions(ctx context.Context, userID string, permissionLevel int32) (*domain.User, error)
	DisableUser(ctx context.Context, userID string) error
	EnableUser(ctx context.Context, userID string) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	ListUsers(ctx context.Context, page, pageSize int32, minPermissionLevel *int32, includeDisabled bool) (*domain.UserList, error)
	DeleteUser(ctx context.Context, userID string) error
}

type SettingsUsecase interface {
	GetSettings(ctx context.Context) (*domain.Settings, error)
	UpdateSettings(ctx context.Context, req *domain.UpdateSettingsRequest) (*domain.Settings, error)
	ResetSettings(ctx context.Context) (*domain.Settings, error)
}
