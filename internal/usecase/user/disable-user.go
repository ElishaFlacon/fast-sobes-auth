package user

import (
	"context"
	"fmt"
)

func (u *usecase) DisableUser(ctx context.Context, userID string) error {
	u.log.Infof("Disable user: %s", userID)

	if err := u.userRepo.SetDisabled(ctx, userID, true); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}
