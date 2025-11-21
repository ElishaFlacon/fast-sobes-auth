package repository

import (
	"context"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

type UserRepository interface {
	// CRUD операции
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error

	// Списки и фильтрация
	List(ctx context.Context, offset, limit int32, minPermissionLevel *int32, includeDisabled bool) ([]*domain.User, int32, error)

	// Управление статусом
	UpdatePermissionLevel(ctx context.Context, id string, level int32) error
	SetDisabled(ctx context.Context, id string, disabled bool) error

	// 2FA операции
	UpdateTwoFactorSecret(ctx context.Context, id, secret string) error
	UpdateTwoFactorEnabled(ctx context.Context, id string, enabled bool) error

	// Изменение email
	UpdateEmail(ctx context.Context, id, newEmail string) error

	// Проверка существования
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type PasswordRepository interface {
	// Хранение хешей паролей
	SetPassword(ctx context.Context, userID, passwordHash string) error
	GetPasswordHash(ctx context.Context, userID string) (string, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
}

type TokenRepository interface {
	// Access tokens
	CreateAccessToken(ctx context.Context, userID string, expiresAt time.Time) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (*domain.TokenInfo, error)
	RevokeAccessToken(ctx context.Context, token string) error

	// Refresh tokens
	CreateRefreshToken(ctx context.Context, userID string, expiresAt time.Time) (string, error)
	VerifyRefreshToken(ctx context.Context, token string) (string, error) // returns userID
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error

	// Временные токены (для регистрации, смены email)
	CreateTempToken(ctx context.Context, userID, tokenType string, data map[string]interface{}, expiresAt time.Time) (string, error)
	VerifyTempToken(ctx context.Context, token, tokenType string) (string, map[string]interface{}, error)
	RevokeTempToken(ctx context.Context, token string) error
}

type TwoFactorRepository interface {
	// Секреты 2FA
	SaveSecret(ctx context.Context, userID, secret string) error
	GetSecret(ctx context.Context, userID string) (string, error)
	DeleteSecret(ctx context.Context, userID string) error

	// Backup коды
	SaveBackupCodes(ctx context.Context, userID string, codes []string) error
	GetBackupCodes(ctx context.Context, userID string) ([]string, error)
	UseBackupCode(ctx context.Context, userID, code string) error
}

type SettingsRepository interface {
	Get(ctx context.Context) (*domain.Settings, error)
	Update(ctx context.Context, settings *domain.Settings) error
	Reset(ctx context.Context) error
}

type EmailRepository interface {
	CreateChangeRequest(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error
	GetChangeRequest(ctx context.Context, token string) (userID, newEmail string, err error)
	DeleteChangeRequest(ctx context.Context, token string) error
	DeleteUserChangeRequests(ctx context.Context, userID string) error
}
