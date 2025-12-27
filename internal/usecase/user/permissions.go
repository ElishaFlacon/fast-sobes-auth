package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

func (u *usecase) UpdatePermissions(ctx context.Context, userID string, permissionLevel int32) (*domain.User, error) {
	u.log.Infof("Update permissions for user id=%s level=%d", userID, permissionLevel)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	user.PermissionLevel = permissionLevel
	user.UpdatedAt = u.now()

	if err := u.users.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return def.SanitizeUser(user), nil
}
