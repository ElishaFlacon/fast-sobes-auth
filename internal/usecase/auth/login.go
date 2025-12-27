package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *usecase) Login(ctx context.Context, email, password string) (*domain.LoginResult, error) {
	u.log.Infof("Login attempt email=%s", email)

	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user.Disabled {
		return nil, fmt.Errorf("user is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	auth, err := u.issueAuth(ctx, user)
	if err != nil {
		return nil, err
	}

	u.log.Infof("User logged in id=%d email=%s", user.ID, user.Email)

	return &domain.LoginResult{Auth: auth}, nil
}
