package user

import (
	"context"
)

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&User{}, &User{Id: id}).Error
}
