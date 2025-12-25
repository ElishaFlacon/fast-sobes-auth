package auth

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// получение пользователя через почту+хеш пароля
// проверка что пользователь не залочен
// выдача токена

func (u *usecase) Login(ctx context.Context, email, password string) (*domain.LoginResult, error) {
	u.log.Infof("Login user: %s", email)

	// Получение пользователя
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Проверка блокировки
	if user.Disabled {
		return nil, fmt.Errorf("user disabled")
	}

	// Проверка пароля
	passwordHash, err := u.passwordRepo.GetPasswordHash(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("get password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	result := &domain.LoginResult{}

	// Создание токенов доступа
	// auth, err := .createAuthTokens(ctx, user)
	// if err != nil {
	// 	return nil, err
	// }
	// result.Auth = auth

	return result, nil
}
