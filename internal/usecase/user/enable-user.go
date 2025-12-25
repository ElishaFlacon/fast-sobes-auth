package user

import (
	"context"
)

func (u *usecase) EnableUser(ctx context.Context, userId string) error {
	u.log.Infof("Enable user: %s", userId)

	return u.userRepo.SetDisabled(ctx, userId, false)
}
