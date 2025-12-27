package domain

import "time"

type AccessToken struct {
	ID        int64
	JTI       string
	UserID    int64
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
