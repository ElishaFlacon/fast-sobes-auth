package app

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	authHandler "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/auth"
	helloHandler "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/hello"
	settingsHandler "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/settings"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	emailRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/email"
	passwordRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/password"
	settingsRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/settings"
	tokenRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/token"
	twoFactorRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/two-factor"
	userRepository "github.com/ElishaFlacon/fast-sobes-auth/internal/repository/user"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	authUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/auth"
	helloUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/hello"
	settingsUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/settings"
	"gorm.io/gorm"
)

type Provider struct {
	cfg *config.Config
	db  *gorm.DB
	log domain.Logger

	// Usecases
	settingsUsecase usecase.SettingsUsecase
	authUsecase     usecase.AuthUsecase
	helloUsecase    usecase.HelloUsecase

	// gRPC handlers
	settingsHandler *settingsHandler.Implementation
	authHandler     *authHandler.Implementation
	helloHandler    *helloHandler.Implementation

	// Repositories
	emailRepository     repository.EmailRepository
	twoFactorRepository repository.TwoFactorRepository
	tokenRepository     repository.TokenRepository
	passwordRepository  repository.PasswordRepository
	userRepository      repository.UserRepository
	settingsRepository  repository.SettingsRepository
}

func NewProvider(cfg *config.Config, db *gorm.DB, log domain.Logger) *Provider {
	return &Provider{
		cfg: cfg,
		db:  db,
		log: log,
	}
}

// ----------------- gRPC HANDLER ------------------

func (p *Provider) SettingsHandler() *settingsHandler.Implementation {
	if p.settingsHandler == nil {
		p.settingsHandler = settingsHandler.NewImplementation(p.SettingsUsecase())
	}
	return p.settingsHandler
}

func (p *Provider) AuthHandler() *authHandler.Implementation {
	if p.authHandler == nil {
		p.authHandler = authHandler.NewImplementation(p.AuthUsecase())
	}
	return p.authHandler
}

func (p *Provider) HelloHandler() *helloHandler.Implementation {
	if p.helloHandler == nil {
		p.helloHandler = helloHandler.NewImplementation(p.HelloUsecase())
	}
	return p.helloHandler
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
			p.PasswordRepository(),
			p.TokenRepository(),
			p.TwoFactorRepository(),
			p.SettingsRepository(),
			p.EmailRepository(),
		)
	}
	return p.authUsecase
}

func (p *Provider) HelloUsecase() usecase.HelloUsecase {
	if p.helloUsecase == nil {
		p.helloUsecase = helloUsecase.NewUsecase(p.log, nil)
	}
	return p.helloUsecase
}

// ------------------ REPOSITORY -------------------

func (p *Provider) EmailRepository() repository.EmailRepository {
	if p.emailRepository == nil {
		p.emailRepository = emailRepository.NewRepository(p.db)
	}
	return p.emailRepository
}

func (p *Provider) TwoFactorRepository() repository.TwoFactorRepository {
	if p.twoFactorRepository == nil {
		p.twoFactorRepository = twoFactorRepository.NewRepository(p.db)
	}
	return p.twoFactorRepository
}

func (p *Provider) TokenRepository() repository.TokenRepository {
	if p.tokenRepository == nil {
		p.tokenRepository = tokenRepository.NewRepository(p.db)
	}
	return p.tokenRepository
}

func (p *Provider) PasswordRepository() repository.PasswordRepository {
	if p.passwordRepository == nil {
		p.passwordRepository = passwordRepository.NewRepository(p.db)
	}
	return p.passwordRepository
}

func (p *Provider) UserRepository() repository.UserRepository {
	if p.userRepository == nil {
		p.userRepository = userRepository.NewRepository(p.db)
	}
	return p.userRepository
}

func (p *Provider) SettingsRepository() repository.SettingsRepository {
	if p.settingsRepository == nil {
		p.settingsRepository = settingsRepository.NewRepository(p.db)
	}
	return p.settingsRepository
}
