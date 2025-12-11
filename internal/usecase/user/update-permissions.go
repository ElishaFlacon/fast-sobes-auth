package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) UpdatePermissions(ctx context.Context, userID string, permissionLevel int32) (*domain.User, error) {
	u.log.Infof("Update permissions for user: %s to level %d", userID, permissionLevel)

	if err := u.userRepo.UpdatePermissionLevel(ctx, userID, permissionLevel); err != nil {
		return nil, fmt.Errorf("update permission level: %w", err)
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
