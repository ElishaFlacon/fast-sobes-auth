package domain

import (
	"time"
)

type User struct {
	Id               int64
	Email            string
	PasswordHash     string
	PermissionLevel  int32
	Disabled         bool
	TwoFactorEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserList struct {
	Users    []*User
	Total    int64
	Page     int32
	PageSize int32
}
