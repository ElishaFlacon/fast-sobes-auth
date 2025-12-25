package domain

import (
	"time"
)

type Settings struct {
	Id                        int64
	RequireTwoFactor          bool
	TokenTTLMinutes           int32
	MinPasswordLength         int32
	RequirePasswordComplexity bool
	UpdatedAt                 time.Time
}
