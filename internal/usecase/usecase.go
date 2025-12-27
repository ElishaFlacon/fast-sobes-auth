package usecase

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

type AuthUsecase interface {
	Register(ctx context.Context, email, password string) (*domain.RegisterResult, error)
	Login(ctx context.Context, email, password string) (*domain.LoginResult, error)
	Logout(ctx context.Context, token string) error
}

type UserUsecase interface {
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	UsersList(
		ctx context.Context,
		page, pageSize int32,
		minPermissionLevel *int32,
		includeDisabled bool,
	) (*domain.UserList, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ChangeEmail(ctx context.Context, userID, newEmail, password string) error
	VerifyEmailChange(ctx context.Context, token string) error
	UpdatePermissions(ctx context.Context, userID string, permissionLevel int32) (*domain.User, error)
	DisableUser(ctx context.Context, userID string) error
	EnableUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
}

type SettingsUsecase interface {
	GetSettings(ctx context.Context) (*domain.Settings, error)
	UpdateSettings(
		ctx context.Context,
		requireTwoFactor *bool,
		tokenTTLMinutes *int32,
		minPasswordLength *int32,
		requirePasswordComplexity *bool,
	) (*domain.Settings, error)
	ResetSettings(ctx context.Context) (*domain.Settings, error)
}

type TokensUsecase interface {
	VerifyToken(ctx context.Context, token string) (*domain.TokenInfo, error)
}

type EmailUsecase interface {
	Send(ctx context.Context, to, subject, body string) error
}
