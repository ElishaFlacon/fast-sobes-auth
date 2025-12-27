package user

import (
	"context"
	"fmt"

	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

func (u *usecase) DisableUser(ctx context.Context, userID string) error {
	u.log.Infof("Disable user id=%s", userID)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	user.Disabled = true
	user.UpdatedAt = u.now()

	if err := u.users.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if err := u.tokens.RevokeAllByUser(ctx, id); err != nil {
		u.log.Errorf("failed to revoke tokens for disabled user %d: %v", id, err)
	}

	return nil
}

func (u *usecase) EnableUser(ctx context.Context, userID string) error {
	u.log.Infof("Enable user id=%s", userID)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	user.Disabled = false
	user.UpdatedAt = u.now()

	if err := u.users.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (u *usecase) DeleteUser(ctx context.Context, userID string) error {
	u.log.Infof("Delete user id=%s", userID)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return err
	}

	if err := u.tokens.RevokeAllByUser(ctx, id); err != nil {
		u.log.Errorf("failed to revoke tokens before delete for user %d: %v", id, err)
	}

	if err := u.users.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
