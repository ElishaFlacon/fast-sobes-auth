package user

import (
	"context"
)

func (u *usecase) EnableUser(ctx context.Context, userID string) error {
	u.log.Infof("Enable user: %s", userID)

	return u.userRepo.SetDisabled(ctx, userID, false)
}
