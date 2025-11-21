package password

import (
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	"gorm.io/gorm"
)

var _ def.PasswordRepository = (*repository)(nil)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}
