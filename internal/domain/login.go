package domain

type LoginResult struct {
	RequiresTwoFactor bool
	TempToken         string
	Auth              *AuthResult
}
