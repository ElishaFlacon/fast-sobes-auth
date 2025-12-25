package domain

import "time"

type AccessToken struct {
	Id        int64
	Token     string
	UserId    int64
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
