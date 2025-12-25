package accessToken

import (
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/repository"
	"gorm.io/gorm"
)

var _ def.AccessTokenRepository = (*repository)(nil)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&AccessToken{},
	)
}
