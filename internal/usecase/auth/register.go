package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) Register(ctx context.Context, email, password string) (*domain.RegisterResult, error) {
	u.log.Infof("Register user: %s", email)

	// Проверка существования пользователя
	exists, err := u.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check user exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Получение настроек
	settings, err := u.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}

	// Валидация пароля
	if err := utils.ValidatePassword(password, settings); err != nil {
		return nil, err
	}

	// Создание пользователя
	user := &domain.User{
		Id:               uuid.New().String(),
		Email:            email,
		PermissionLevel:  0,
		Disabled:         false,
		TwoFactorEnabled: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	if err := u.passwordRepo.SetPassword(ctx, user.Id, string(hash)); err != nil {
		return nil, fmt.Errorf("save password: %w", err)
	}

	result := &domain.RegisterResult{
		User:              user,
		RequiresTwoFactor: settings.RequireTwoFactor,
	}

	return result, nil
}
