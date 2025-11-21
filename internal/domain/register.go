package domain

type RegisterResult struct {
	User              *User
	TwoFactorSecret   string
	QRCodeURL         string
	RequiresTwoFactor bool
	TempToken         string
}
