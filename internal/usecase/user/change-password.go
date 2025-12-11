package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	u.log.Infof("Change password for user: %s", userID)

	// Проверка старого пароля
	oldHash, err := u.passwordRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Валидация нового пароля
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("get settings: %w", err)
	}

	if err := utils.ValidatePassword(newPassword, settings); err != nil {
		return err
	}

	// Хеширование нового пароля
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Обновление пароля
	if err := u.passwordRepo.UpdatePassword(ctx, userID, string(newHash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}