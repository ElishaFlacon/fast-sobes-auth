package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	u.log.Infof("Get user: %s", userId)

	return u.userRepo.GetById(ctx, userId)
}
