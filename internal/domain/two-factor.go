package domain

type TwoFactorSetup struct {
	Secret      string
	QRCodeURL   string
	BackupCodes []string
}
