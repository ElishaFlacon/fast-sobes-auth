package auth

import (
	"context"
	"fmt"
	"time"

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

	if err := u.validatePassword(newPassword, settings); err != nil {
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

func (u *usecase) ChangeEmail(ctx context.Context, userID, newEmail, password string) error {
	u.log.Infof("Change email for user: %s", userID)

	// Проверка пароля
	passwordHash, err := u.passwordRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	// Проверка существования email
	exists, err := u.userRepo.ExistsByEmail(ctx, newEmail)
	if err != nil {
		return fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return fmt.Errorf("email already in use")
	}

	// Удаление старых запросов на смену email
	if err := u.emailRepo.DeleteUserChangeRequests(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete old change requests: %v", err)
	}

	// Создание токена для подтверждения
	token, err := u.tokenRepo.CreateTempToken(
		ctx,
		userID,
		"email_change",
		map[string]interface{}{"new_email": newEmail},
		time.Now().Add(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("create verification token: %w", err)
	}

	// Сохранение запроса на смену email
	if err := u.emailRepo.CreateChangeRequest(
		ctx,
		userID,
		newEmail,
		token,
		time.Now().Add(24*time.Hour),
	); err != nil {
		return fmt.Errorf("create change request: %w", err)
	}

	// TODO: отправить email с подтверждением

	return nil
}

func (u *usecase) VerifyEmailChange(ctx context.Context, token string) error {
	u.log.Infof("Verify email change")

	// Получение запроса на смену email
	userID, newEmail, err := u.emailRepo.GetChangeRequest(ctx, token)
	if err != nil {
		return fmt.Errorf("get change request: %w", err)
	}

	// Обновление email
	if err := u.userRepo.UpdateEmail(ctx, userID, newEmail); err != nil {
		return fmt.Errorf("update email: %w", err)
	}

	// Удаление запроса
	if err := u.emailRepo.DeleteChangeRequest(ctx, token); err != nil {
		u.log.Errorf("Failed to delete change request: %v", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}
