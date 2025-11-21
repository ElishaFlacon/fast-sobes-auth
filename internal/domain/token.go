package domain

type TokenInfo struct {
	Valid           bool
	UserID          string
	PermissionLevel int32
	Email           string
}
