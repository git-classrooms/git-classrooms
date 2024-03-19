package session

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
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

		// Wir können typen registrieren, die in der Session gespeichert werden sollen
		// Damit sparen wir uns das aufteilen der structs in einzelne Felder
		store.RegisterType(time.Time{})
		store.RegisterType(UserState(0))
		store.RegisterType(&oauth2.Token{})
	})
}
