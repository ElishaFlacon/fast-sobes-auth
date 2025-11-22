package twofactor

import (
	"context"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func (r *repository) SaveSecret(ctx context.Context, userID, secret string) error {
	twoFactorSecret := &models.TwoFactorSecret{
		UserID: userID,
		Secret: secret,
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"secret", "updated_at"}),
		}).
		Create(twoFactorSecret).Error
}

func (r *repository) GetSecret(ctx context.Context, userID string) (string, error) {
	var secret models.TwoFactorSecret

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&secret).Error; err != nil {
		return "", err
	}

	return secret.Secret, nil
}

func (r *repository) DeleteSecret(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.TwoFactorSecret{}, "user_id = ?", userID).Error
}

func (r *repository) SaveBackupCodes(ctx context.Context, userID string, codes []string) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Удаление старых кодов
	if err := tx.Where("user_id = ?", userID).Delete(&models.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Хеширование и сохранение новых кодов
	for _, code := range codes {
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return err
		}

		backupCode := &models.BackupCode{
			UserID: userID,
			Code:   string(hash),
			Used:   false,
		}

		if err := tx.Create(backupCode).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) GetBackupCodes(ctx context.Context, userID string) ([]string, error) {
	var codes []models.BackupCode

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND used = ?", userID, false).
		Find(&codes).Error; err != nil {
		return nil, err
	}

	result := make([]string, len(codes))
	for i, code := range codes {
		result[i] = code.Code
	}

	return result, nil
}

func (r *repository) UseBackupCode(ctx context.Context, userID, code string) error {
	var codes []models.BackupCode

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND used = ?", userID, false).
		Find(&codes).Error; err != nil {
		return err
	}

	// Проверка кода
	for _, bc := range codes {
		if err := bcrypt.CompareHashAndPassword([]byte(bc.Code), []byte(code)); err == nil {
			// Найден правильный код, помечаем как использованный
			now := time.Now()
			return r.db.WithContext(ctx).Model(&models.BackupCode{}).
				Where("id = ?", bc.ID).
				Updates(map[string]interface{}{
					"used":    true,
					"used_at": &now,
				}).Error
		}
	}

	return fmt.Errorf("invalid backup code")
}
