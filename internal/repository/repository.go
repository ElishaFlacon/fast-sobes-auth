package repository

import (
	"context"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

type UserRepository interface {
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetList(
		ctx context.Context,
		offset,
		limit int32,
		minPermissionLevel *int32,
		includeDisabled bool,
	) ([]*domain.User, int64, error)
	Create(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

type AccessTokenRepository interface {
	Create(ctx context.Context, token *domain.AccessToken) error
	GetByToken(ctx context.Context, token string) (*domain.AccessToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllByUser(ctx context.Context, userId int64) error
	DeleteExpired(ctx context.Context, now time.Time) error
}
type SettingsRepository interface {
	Get(ctx context.Context) (*domain.Settings, error)
	Update(ctx context.Context, settings *domain.Settings) error
	Reset(ctx context.Context) error
}
