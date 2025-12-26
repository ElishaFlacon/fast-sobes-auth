package auth

import (
	"os"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	emailStub "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/email"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/jwtmanager"
)

var _ def.AuthUsecase = (*usecase)(nil)

type usecase struct {
	log      domain.Logger
	users    repository.UserRepository
	tokens   repository.AccessTokenRepository
	settings repository.SettingsRepository
	email    def.EmailUsecase
	jwt      *jwtmanager.Manager
	now      func() time.Time
}

func NewUsecase(
	log domain.Logger,
	users repository.UserRepository,
	tokens repository.AccessTokenRepository,
	settings repository.SettingsRepository,
	email def.EmailUsecase,
	jwtSecret string,
) *usecase {
	if email == nil {
		email = emailStub.NewUsecase(log)
	}

	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}

	return &usecase{
		log:      log,
		users:    users,
		tokens:   tokens,
		settings: settings,
		email:    email,
		jwt:      jwtmanager.New(jwtSecret),
		now:      time.Now,
	}
}
