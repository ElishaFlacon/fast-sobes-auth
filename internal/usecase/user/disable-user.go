package user

import (
	"context"
	"fmt"
)

func (u *usecase) DisableUser(ctx context.Context, userId string) error {
	u.log.Infof("Disable user: %s", userId)

	if err := u.userRepo.SetDisabled(ctx, userId, true); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userId); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}
