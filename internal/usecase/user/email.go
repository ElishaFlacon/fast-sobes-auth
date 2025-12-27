package user

import (
	"context"
	"errors"
	"fmt"

	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *usecase) ChangeEmail(ctx context.Context, userID, newEmail, password string) error {
	u.log.Infof("Change email for user id=%s new_email=%s", userID, newEmail)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	if user.Email == newEmail {
		return fmt.Errorf("new email matches current email")
	}

	_, err = u.users.GetByEmail(ctx, newEmail)
	if err == nil {
		return fmt.Errorf("email %s is already used", newEmail)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check email availability: %w", err)
	}

	user.Email = newEmail
	user.UpdatedAt = u.now()

	if err := u.users.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if err := u.tokens.RevokeAllByUser(ctx, user.ID); err != nil {
		u.log.Errorf("failed to revoke tokens after email change: %v", err)
	}

	_ = u.email.Send(ctx, newEmail, "Email changed", fmt.Sprintf("Email for user %d updated", user.ID))

	u.log.Infof("Email updated for user id=%d", user.ID)

	return nil
}

func (u *usecase) VerifyEmailChange(ctx context.Context, token string) error {
	u.log.Infof("Email verification stub token=%s", token)
	_ = u.email.Send(ctx, "", "Verify email stub", fmt.Sprintf("Received verification token %s", token))
	return nil
}
