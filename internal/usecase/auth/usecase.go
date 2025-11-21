package auth

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

var _ def.AuthUsecase = (*usecase)(nil)

type usecase struct {
	log           domain.Logger
	userRepo      repository.UserRepository
	passwordRepo  repository.PasswordRepository
	tokenRepo     repository.TokenRepository
	twoFactorRepo repository.TwoFactorRepository
	settingsRepo  repository.SettingsRepository
	emailRepo     repository.EmailRepository
}

func NewUsecase(
	log domain.Logger,
	userRepo repository.UserRepository,
	passwordRepo repository.PasswordRepository,
	tokenRepo repository.TokenRepository,
	twoFactorRepo repository.TwoFactorRepository,
	settingsRepo repository.SettingsRepository,
	emailRepo repository.EmailRepository,
) *usecase {
	return &usecase{
		log:           log,
		userRepo:      userRepo,
		passwordRepo:  passwordRepo,
		tokenRepo:     tokenRepo,
		twoFactorRepo: twoFactorRepo,
		settingsRepo:  settingsRepo,
		emailRepo:     emailRepo,
	}
}
