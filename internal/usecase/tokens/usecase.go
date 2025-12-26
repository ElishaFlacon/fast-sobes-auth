package tokens

import (
	"os"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase/jwtmanager"
)

var _ def.TokensUsecase = (*usecase)(nil)

type usecase struct {
	log    domain.Logger
	tokens repository.AccessTokenRepository
	users  repository.UserRepository
	jwt    *jwtmanager.Manager
	now    func() time.Time
}

func NewUsecase(
	log domain.Logger,
	tokens repository.AccessTokenRepository,
	users repository.UserRepository,
	jwtSecret string,
) *usecase {
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}

	return &usecase{
		log:    log,
		tokens: tokens,
		users:  users,
		jwt:    jwtmanager.New(jwtSecret),
		now:    time.Now,
	}
}
