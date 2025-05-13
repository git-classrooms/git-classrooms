package stores

import (
	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"gorm.io/gorm"
)

func NewSessionStore(db *gorm.DB) (scs.Store, error) {
	return gormstore.New(db)
}
