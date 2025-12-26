package auth

import (
	"context"
	"fmt"
)

func (u *usecase) Logout(ctx context.Context, token string) error {
	u.log.Infof("Logout")

	if token == "" {
		return fmt.Errorf("token is required")
	}

	if err := u.tokens.Revoke(ctx, token); err != nil {
		return fmt.Errorf("revoke access token: %w", err)
	}

	return nil
}
