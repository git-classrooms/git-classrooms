package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"sync"
	"time"
)

var store *session.Store
var once sync.Once

func InitSessionStore(dsn string) {
	once.Do(func() {
		store = session.New(session.Config{
			Expiration:     time.Hour * 24 * 7,
			CookieHTTPOnly: true,
			CookieSecure:   true,
			//Storage: postgres.New(postgres.Config{
			//	ConnectionURI: dsn,
			//	Table:         "fiber_storage",
			//	Reset:         false,
			//	GCInterval:    10 * time.Second,
			//}),
		})
	})
}
