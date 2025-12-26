package app

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	accessTokenRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/access-token"
	settingsRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/settings"
	userRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/user"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	authUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/auth"
	emailUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/email"
	settingsUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/settings"
	tokensUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/tokens"
	userUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/user"
	"gorm.io/gorm"
)

type Provider struct {
	cfg *config.Config
	db  *gorm.DB
	log domain.Logger

	// Usecases
	settingsUsecase usecase.SettingsUsecase
	authUsecase     usecase.AuthUsecase
	userUsecase     usecase.UserUsecase
	tokensUsecase   usecase.TokensUsecase
	emailUsecase    usecase.EmailUsecase

	// Repositories
	userRepository        repository.UserRepository
	accessTokenRepository repository.AccessTokenRepository
	settingsRepository    repository.SettingsRepository
}

func NewProvider(cfg *config.Config, db *gorm.DB, log domain.Logger) *Provider {
	return &Provider{
		cfg: cfg,
		db:  db,
		log: log,
	}
}

// -------------------- USECASE --------------------

func (p *Provider) SettingsUsecase() usecase.SettingsUsecase {
	if p.settingsUsecase == nil {
		p.settingsUsecase = settingsUsecase.NewUsecase(p.log, p.SettingsRepository())
	}
	return p.settingsUsecase
}

func (p *Provider) AuthUsecase() usecase.AuthUsecase {
	if p.authUsecase == nil {
		p.authUsecase = authUsecase.NewUsecase(
			p.log,
			p.UserRepository(),
			p.AccessTokenRepository(),
			p.SettingsRepository(),
			p.EmailUsecase(),
			"",
		)
	}
	return p.authUsecase
}

func (p *Provider) UserUsecase() usecase.UserUsecase {
	if p.userUsecase == nil {
		p.userUsecase = userUsecase.NewUsecase(
			p.log,
			p.UserRepository(),
			p.SettingsRepository(),
			p.AccessTokenRepository(),
			p.EmailUsecase(),
		)
	}
	return p.userUsecase
}

func (p *Provider) TokensUsecase() usecase.TokensUsecase {
	if p.tokensUsecase == nil {
		p.tokensUsecase = tokensUsecase.NewUsecase(
			p.log,
			p.AccessTokenRepository(),
			p.UserRepository(),
			"",
		)
	}
	return p.tokensUsecase
}

func (p *Provider) EmailUsecase() usecase.EmailUsecase {
	if p.emailUsecase == nil {
		p.emailUsecase = emailUsecase.NewUsecase(p.log)
	}
	return p.emailUsecase
}

// ------------------ REPOSITORY -------------------

func (p *Provider) UserRepository() repository.UserRepository {
	if p.userRepository == nil {
		p.userRepository = userRepository.NewRepository(p.db)
	}
	return p.userRepository
}

func (p *Provider) AccessTokenRepository() repository.AccessTokenRepository {
	if p.accessTokenRepository == nil {
		p.accessTokenRepository = accessTokenRepository.NewRepository(p.db)
	}
	return p.accessTokenRepository
}

func (p *Provider) SettingsRepository() repository.SettingsRepository {
	if p.settingsRepository == nil {
		p.settingsRepository = settingsRepository.NewRepository(p.db)
	}
	return p.settingsRepository
}
