package domain

import (
	"time"
)

type Settings struct {
	ID                        string
	RequireTwoFactor          bool
	TokenTTLMinutes           int32
	RefreshTokenTTLDays       int32
	MinPasswordLength         int32
	RequirePasswordComplexity bool
	UpdatedAt                 time.Time
}
