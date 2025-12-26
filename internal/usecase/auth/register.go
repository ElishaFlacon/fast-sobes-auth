package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *usecase) Register(ctx context.Context, email, password string) (*domain.RegisterResult, error) {
	u.log.Infof("Register user email=%s", email)

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	_, err := u.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	settings, err := u.loadSettings(ctx)
	if err != nil {
		return nil, err
	}

	if err := def.ValidatePassword(password, settings); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := u.now()
	user := &domain.User{
		Email:           email,
		PasswordHash:    string(hash),
		PermissionLevel: 1,
		Disabled:        false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	u.log.Infof("User registered id=%d email=%s", user.Id, user.Email)

	return &domain.RegisterResult{
		User: def.SanitizeUser(user),
	}, nil
}
