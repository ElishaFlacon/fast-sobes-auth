package user

import (
	"context"
	"fmt"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (u *usecase) UsersList(
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
