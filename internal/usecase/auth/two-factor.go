package auth

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) EnableTwoFactor(ctx context.Context, userID, password string) (*domain.TwoFactorSetup, error) {
	u.log.Infof("Enable 2FA for user: %s", userID)

	// Проверка пароля
	passwordHash, err := u.passwordRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Получение пользователя
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Генерация секрета
	secret, qrCode, err := u.generate2FASecret(user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate 2FA secret: %w", err)
	}

	// Генерация backup кодов
	backupCodes := u.generateBackupCodes()

	// Сохранение секрета и backup кодов
	if err := u.twoFactorRepo.SaveSecret(ctx, userID, secret); err != nil {
		return nil, fmt.Errorf("save 2FA secret: %w", err)
	}

	if err := u.twoFactorRepo.SaveBackupCodes(ctx, userID, backupCodes); err != nil {
		return nil, fmt.Errorf("save backup codes: %w", err)
	}

	// Включение 2FA
	if err := u.userRepo.UpdateTwoFactorEnabled(ctx, userID, true); err != nil {
		return nil, fmt.Errorf("enable 2FA: %w", err)
	}

	return &domain.TwoFactorSetup{
		Secret:      secret,
		QRCodeURL:   qrCode,
		BackupCodes: backupCodes,
	}, nil
}

func (u *usecase) DisableTwoFactor(ctx context.Context, userID, password, code string) error {
	u.log.Infof("Disable 2FA for user: %s", userID)

	// Проверка пароля
	passwordHash, err := u.passwordRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	// Проверка кода 2FA
	secret, err := u.twoFactorRepo.GetSecret(ctx, userID)
	if err != nil {
		return fmt.Errorf("get 2FA secret: %w", err)
	}

	if !u.verify2FACode(secret, code) {
		return fmt.Errorf("invalid 2FA code")
	}

	// Отключение 2FA
	if err := u.userRepo.UpdateTwoFactorEnabled(ctx, userID, false); err != nil {
		return fmt.Errorf("disable 2FA: %w", err)
	}

	// Удаление секрета и backup кодов
	if err := u.twoFactorRepo.DeleteSecret(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete 2FA secret: %v", err)
	}

	return nil
}

func (u *usecase) generateBackupCodes() []string {
	// TODO: реализовать генерацию backup кодов
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		codes[i] = fmt.Sprintf("BACKUP-%d", i)
	}
	return codes
}
