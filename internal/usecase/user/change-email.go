package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
