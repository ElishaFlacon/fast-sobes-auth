package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
)

func (r *repository) CreateAccessToken(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	token := r.generateToken()

	accessToken := &models.AccessToken{
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
	var accessToken models.AccessToken

	err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&accessToken).Error
	if err != nil {
		return nil, err
	}

	// Получение информации о пользователе
	var user models.User
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
	return r.db.WithContext(ctx).Model(&models.AccessToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *repository) CreateRefreshToken(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	token := r.generateToken()

	refreshToken := &models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := r.db.WithContext(ctx).Create(refreshToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (r *repository) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	var refreshToken models.RefreshToken

	err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&refreshToken).Error
	if err != nil {
		return "", err
	}

	return refreshToken.UserID, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
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
	if err := tx.Model(&models.AccessToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Отзыв refresh tokens
	if err := tx.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Отзыв temp tokens
	if err := tx.Model(&models.TempToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) CreateTempToken(
	ctx context.Context, userID,
	tokenType string,
	data map[string]interface{},
	expiresAt time.Time,
) (string, error) {
	token := r.generateToken()

	tempToken := &models.TempToken{
		Token:     token,
		UserID:    userID,
		TokenType: tokenType,
		Data:      data,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := r.db.WithContext(ctx).Create(tempToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (r *repository) VerifyTempToken(ctx context.Context, token, tokenType string) (string, map[string]interface{}, error) {
	var tempToken models.TempToken

	err := r.db.WithContext(ctx).
		Where("token = ? AND token_type = ? AND revoked = ? AND expires_at > ?",
			token, tokenType, false, time.Now()).
		First(&tempToken).Error
	if err != nil {
		return "", nil, err
	}

	return tempToken.UserID, tempToken.Data, nil
}

func (r *repository) RevokeTempToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.TempToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *repository) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
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

	if err := tx.Where("expires_at < ?", now).Delete(&models.AccessToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("expires_at < ?", now).Delete(&models.RefreshToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("expires_at < ?", now).Delete(&models.TempToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
