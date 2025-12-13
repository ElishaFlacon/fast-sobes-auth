package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func (r *repository) List(
	ctx context.Context,
	offset,
	limit int32,
	minPermissionLevel *int32,
	includeDisabled bool,
) ([]*domain.User, int32, error) {
	var modelUsers []User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	if minPermissionLevel != nil {
		query = query.Where("permission_level >= ?", *minPermissionLevel)
	}

	if !includeDisabled {
		query = query.Where("disabled = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&modelUsers).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(modelUsers))
	for i, model := range modelUsers {
		users[i] = r.toDomain(&model)
	}

	return users, int32(total), nil
}
