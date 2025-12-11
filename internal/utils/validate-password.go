package utils

import (
	"fmt"
	"unicode"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func ValidatePassword(password string, settings *domain.Settings) error {
	if len(password) < int(settings.MinPasswordLength) {
		return fmt.Errorf("password too short: minimum %d characters required", settings.MinPasswordLength)
	}

	if settings.RequirePasswordComplexity {
		var (
			hasUpper   bool
			hasLower   bool
			hasNumber  bool
			hasSpecial bool
		)

		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		if !hasUpper {
			return fmt.Errorf("password must contain at least one uppercase letter")
		}
		if !hasLower {
			return fmt.Errorf("password must contain at least one lowercase letter")
		}
		if !hasNumber {
			return fmt.Errorf("password must contain at least one digit")
		}
		if !hasSpecial {
			return fmt.Errorf("password must contain at least one special character")
		}
	}

	return nil
}
