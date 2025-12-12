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

type AccessTokenRepository interface {
	CreateAccessToken(ctx context.Context, userID string, expiresAt time.Time) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (*domain.TokenInfo, error)
	RevokeAccessToken(ctx context.Context, token string) error
}

type SettingsRepository interface {
	Get(ctx context.Context) (*domain.Settings, error)
	Update(ctx context.Context, settings *domain.Settings) error
	Reset(ctx context.Context) error
}
