package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) UpdatePermissions(ctx context.Context, userId string, permissionLevel int32) (*domain.User, error) {
	u.log.Infof("Update permissions for user: %s to level %d", userId, permissionLevel)

	if err := u.userRepo.UpdatePermissionLevel(ctx, userId, permissionLevel); err != nil {
		return nil, fmt.Errorf("update permission level: %w", err)
	}

	user, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
