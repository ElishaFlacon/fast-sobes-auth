package domain

import "time"

type AuthResult struct {
	AccessToken string
	ExpiresAt   time.Time
	User        *User
}

type LoginResult struct {
	Auth *AuthResult
}

type RegisterResult struct {
	User *User
}
