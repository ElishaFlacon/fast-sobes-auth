package accessToken

import (
	"context"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) CreateAccessToken(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	token := r.generateToken()

	accessToken := &AccessToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := r.db.WithContext(ctx).Create(accessToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (r *repository) VerifyAccessToken(ctx context.Context, token string) (*domain.TokenInfo, error) {
	var accessToken AccessToken

	err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&accessToken).Error
	if err != nil {
		return nil, err
	}

	// Получение информации о пользователе
	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", accessToken.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	return &domain.TokenInfo{
		Valid:           true,
		UserID:          user.ID,
		PermissionLevel: user.PermissionLevel,
		Email:           user.Email,
	}, nil
}

func (r *repository) RevokeAccessToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&AccessToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Отзыв access tokens
	if err := tx.Model(&AccessToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Очистка истекших токенов (можно запускать по расписанию)
func (r *repository) CleanupExpiredTokens(ctx context.Context) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	if err := tx.Where("expires_at < ?", now).Delete(&AccessToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
