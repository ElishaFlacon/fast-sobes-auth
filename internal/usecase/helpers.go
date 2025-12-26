package usecase

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func SanitizeUser(user *domain.User) *domain.User {
	if user == nil {
		return nil
	}

	clean := *user
	clean.PasswordHash = ""

	return &clean
}

func ParseUserID(id string) (int64, error) {
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: %w", err)
	}

	return userID, nil
}

func ValidatePassword(password string, settings *domain.Settings) error {
	if settings == nil {
		return fmt.Errorf("password settings not provided")
	}

	if int32(len(password)) < settings.MinPasswordLength {
		return fmt.Errorf("password is too short; need at least %d characters", settings.MinPasswordLength)
	}

	if settings.RequirePasswordComplexity {
		var hasLetter bool
		var hasDigit bool

		for _, r := range password {
			if unicode.IsLetter(r) {
				hasLetter = true
			}
			if unicode.IsDigit(r) {
				hasDigit = true
			}
		}

		if !hasLetter || !hasDigit {
			return fmt.Errorf("password must contain both letters and digits")
		}
	}

	return nil
}

func DefaultSettings() *domain.Settings {
	return &domain.Settings{
		RequireTwoFactor:          false,
		TokenTTLMinutes:           60,
		MinPasswordLength:         8,
		RequirePasswordComplexity: true,
	}
}
