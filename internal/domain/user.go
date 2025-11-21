package domain

import (
	"time"
)

type User struct {
	ID               string
	Email            string
	PermissionLevel  int32
	Disabled         bool
	TwoFactorEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserList struct {
	Users    []*User
	Total    int32
	Page     int32
	PageSize int32
}
