package user

import (
	"context"
	"fmt"
)

func (u *usecase) DeleteUser(ctx context.Context, userId string) error {
	u.log.Infof("Delete user: %s", userId)

	// Удаление всех связанных данных
	if err := u.twoFactorRepo.DeleteSecret(ctx, userId); err != nil {
		u.log.Errorf("Failed to delete 2FA secret: %v", err)
	}

	if err := u.emailRepo.DeleteUserChangeRequests(ctx, userId); err != nil {
		u.log.Errorf("Failed to delete change requests: %v", err)
	}

	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userId); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	// Удаление пользователя
	if err := u.userRepo.Delete(ctx, userId); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
