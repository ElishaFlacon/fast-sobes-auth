package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	u.log.Infof("Get user: %s", userID)

	return u.userRepo.GetByID(ctx, userID)
}
