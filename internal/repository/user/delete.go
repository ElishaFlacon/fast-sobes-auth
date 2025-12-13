package user

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)



func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}