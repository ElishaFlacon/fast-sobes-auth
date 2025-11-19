package app

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	helloHandler "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/hello"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	helloUsecase "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/hello"
)

type Provider struct {
	cfg *config.Config
	log domain.Logger

	// Usecases
	helloUsecase usecase.HelloUsecase

	// gRPC handlers
	helloHandler *helloHandler.Implementation

	// Repositories
	// ---
}

func NewProvider(cfg *config.Config, log domain.Logger) *Provider {
	return &Provider{
		cfg: cfg,
		log: log,
	}
}

// ----------------- gRPC HANDLER ------------------

func (p *Provider) HelloHandler() *helloHandler.Implementation {
	if p.helloHandler == nil {
		p.helloHandler = helloHandler.NewImplementation(p.HelloUsecase())
	}
	return p.helloHandler
}

// -------------------- SERVICE --------------------

func (p *Provider) HelloUsecase() usecase.HelloUsecase {
	if p.helloUsecase == nil {
		p.helloUsecase = helloUsecase.NewUsecase(p.log, nil)
	}
	return p.helloUsecase
}

// ------------------ REPOSITORY -------------------
