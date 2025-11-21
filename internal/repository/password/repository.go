package password

import (
    def "github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
)

var _ def.PasswordRepository = (*repository)(nil)

type repository struct {
}

func NewRepository() *repository {
    return &repository{}
}
