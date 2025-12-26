package user

import (
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	emailStub "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/email"
)

var _ def.UserUsecase = (*usecase)(nil)

type usecase struct {
	log      domain.Logger
	users    repository.UserRepository
	settings repository.SettingsRepository
	tokens   repository.AccessTokenRepository
	email    def.EmailUsecase
	now      func() time.Time
}

func NewUsecase(
	log domain.Logger,
	users repository.UserRepository,
	settings repository.SettingsRepository,
	tokens repository.AccessTokenRepository,
	email def.EmailUsecase,
) *usecase {
	if email == nil {
		email = emailStub.NewUsecase(log)
	}

	return &usecase{
		log:      log,
		users:    users,
		settings: settings,
		tokens:   tokens,
		email:    email,
		now:      time.Now,
	}
}
