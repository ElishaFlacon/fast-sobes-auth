package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

func (u *usecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	u.log.Infof("Get user id=%s", userID)

	id, err := def.ParseUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return def.SanitizeUser(user), nil
}

func (u *usecase) UsersList(
	ctx context.Context,
	page, pageSize int32,
	minPermissionLevel *int32,
	includeDisabled bool,
) (*domain.UserList, error) {
	u.log.Infof("List users page=%d pageSize=%d includeDisabled=%t", page, pageSize, includeDisabled)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	users, total, err := u.users.GetList(ctx, offset, pageSize, minPermissionLevel, includeDisabled)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	for i, usr := range users {
		users[i] = def.SanitizeUser(usr)
	}

	return &domain.UserList{
		Users:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
