package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"sync"
)

var store *session.Store
var once sync.Once

func init() {
	once.Do(func() {
		store = session.New()
	})
}
