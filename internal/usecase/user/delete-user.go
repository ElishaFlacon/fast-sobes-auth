package user

import (
	"context"
	"fmt"
)

func (u *usecase) DeleteUser(ctx context.Context, userID string) error {
	u.log.Infof("Delete user: %s", userID)

	// Удаление всех связанных данных
	if err := u.twoFactorRepo.DeleteSecret(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete 2FA secret: %v", err)
	}

	if err := u.emailRepo.DeleteUserChangeRequests(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete change requests: %v", err)
	}

	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	// Удаление пользователя
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
