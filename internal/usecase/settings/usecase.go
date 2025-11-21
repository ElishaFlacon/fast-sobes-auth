package settings

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

var _ def.SettingsUsecase = (*usecase)(nil)

type usecase struct {
	log          domain.Logger
	settingsRepo repository.SettingsRepository
}

func NewUsecase(
	log domain.Logger,
	settingsRepo repository.SettingsRepository,
) *usecase {
	return &usecase{
		log:          log,
		settingsRepo: settingsRepo,
	}
}
