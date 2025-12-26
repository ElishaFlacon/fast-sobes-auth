package email

import (
	"context"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
)

var _ def.EmailUsecase = (*usecase)(nil)

type usecase struct {
	log domain.Logger
}

func NewUsecase(log domain.Logger) *usecase {
	return &usecase{log: log}
}

func (u *usecase) Send(_ context.Context, to, subject, body string) error {
	u.log.Infof("Email stub: to=%s subject=%s body=%s", to, subject, body)
	return nil
}
