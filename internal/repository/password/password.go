package password

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
)

func (r *repository) SetPassword(ctx context.Context, userID, passwordHash string) error {
	password := &models.Password{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	return r.db.WithContext(ctx).Create(password).Error
}

func (r *repository) GetPasswordHash(ctx context.Context, userID string) (string, error) {
	var password models.Password

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&password).Error; err != nil {
		return "", err
	}

	return password.PasswordHash, nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	return r.db.WithContext(ctx).Model(&models.Password{}).
		Where("user_id = ?", userID).
		Update("password_hash", newPasswordHash).Error
}
