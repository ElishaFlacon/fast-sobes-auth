package user

import (
	"context"
	"fmt"

	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	u.log.Infof("Change password for user id=%s", userID)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	settings, err := u.loadSettings(ctx)
	if err != nil {
		return err
	}

	if err := def.ValidatePassword(newPassword, settings); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = string(hash)
	user.UpdatedAt = u.now()

	if err := u.users.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if err := u.tokens.RevokeAllByUser(ctx, user.ID); err != nil {
		u.log.Errorf("failed to revoke tokens after password change: %v", err)
	}

	u.log.Infof("Password changed for user id=%d", user.ID)

	return nil
}
