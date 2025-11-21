package email

import (
	"context"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
)

func (r *repository) CreateChangeRequest(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error {
	request := &models.EmailChangeRequest{
		UserID:    userID,
		NewEmail:  newEmail,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return r.db.WithContext(ctx).Create(request).Error
}

func (r *repository) GetChangeRequest(ctx context.Context, token string) (userID, newEmail string, err error) {
	var request models.EmailChangeRequest

	err = r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&request).Error
	if err != nil {
		return "", "", err
	}

	return request.UserID, request.NewEmail, nil
}

func (r *repository) DeleteChangeRequest(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&models.EmailChangeRequest{}, "token = ?", token).Error
}

func (r *repository) DeleteUserChangeRequests(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.EmailChangeRequest{}, "user_id = ?", userID).Error
}

// Очистка истекших запросов
func (r *repository) CleanupExpiredRequests(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.EmailChangeRequest{}).Error
}
