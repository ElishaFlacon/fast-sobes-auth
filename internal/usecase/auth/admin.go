package auth

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

func (u *usecase) DisableUser(ctx context.Context, userID string) error {
	u.log.Infof("Disable user: %s", userID)

	if err := u.userRepo.SetDisabled(ctx, userID, true); err != nil {
		return fmt.Errorf("disable user: %w", err)
	}

	// Отзыв всех токенов пользователя
	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	return nil
}

func (u *usecase) EnableUser(ctx context.Context, userID string) error {
	u.log.Infof("Enable user: %s", userID)

	return u.userRepo.SetDisabled(ctx, userID, false)
}

func (u *usecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	u.log.Infof("Get user: %s", userID)

	return u.userRepo.GetByID(ctx, userID)
}

func (u *usecase) ListUsers(
	ctx context.Context,
	page,
	pageSize int32,
	minPermissionLevel *int32,
	includeDisabled bool,
) (*domain.UserList, error) {
	u.log.Infof("List users: page=%d, pageSize=%d", page, pageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	users, total, err := u.userRepo.List(ctx, offset, pageSize, minPermissionLevel, includeDisabled)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return &domain.UserList{
		Users:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (u *usecase) DeleteUser(ctx context.Context, userID string) error {
	u.log.Infof("Delete user: %s", userID)

	// Удаление всех связанных данных
	if err := u.twoFactorRepo.DeleteSecret(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete 2FA secret: %v", err)
	}

	if err := u.emailRepo.DeleteUserChangeRequests(ctx, userID); err != nil {
		u.log.Errorf("Failed to delete change requests: %v", err)
	}

	if err := u.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		u.log.Errorf("Failed to revoke user tokens: %v", err)
	}

	// Удаление пользователя
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
