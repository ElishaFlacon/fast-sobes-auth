package user

import (
	"context"
	"fmt"
)

func (u *usecase) VerifyEmailChange(ctx context.Context, token string) error {
	u.log.Infof("Verify email change")

	// Получение запроса на смену email
	userId, newEmail, err := u.emailRepo.GetChangeRequest(ctx, token)
	if err != nil {
		return fmt.Errorf("get change request: %w", err)
	}

	// Обновление email
	if err := u.userRepo.UpdateEmail(ctx, userId, newEmail); err != nil {
		return fmt.Errorf("update email: %w", err)
	}

	// Удаление запроса
	if err := u.emailRepo.DeleteChangeRequest(ctx, token); err != nil {
		u.log.Errorf("Failed to delete change request: %v", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userId); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}
