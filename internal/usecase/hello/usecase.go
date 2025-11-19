package hello

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

var _ def.HelloUsecase = (*usecase)(nil)

type usecase struct {
	log        domain.Logger
	repository *struct{}
}

func NewUsecase(
	log domain.Logger,
	repository *struct{},
) *usecase {
	return &usecase{
		log,
		repository,
	}
}
