package session

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/oauth2"
)

const (
	HeaderName  = "X-Csrf-Token"
	FormKeyName = "csrf_token"
)

var (
	formExtractor   = csrf.CsrfFromForm(FormKeyName)
	headerExtractor = csrf.CsrfFromHeader(HeaderName)
)

var store *session.Store
var CsrfConfig csrf.Config
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

		CsrfConfig = csrf.Config{
			Next: func(c *fiber.Ctx) bool {
				return false
			},
			CookieName:        "csrf_", // __Host-csrf_
			CookieSameSite:    "Lax",
			CookieSecure:      true,
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			Expiration:        1 * time.Hour,
			KeyGenerator:      utils.UUIDv4,
			SingleUseToken:    true,
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return fiber.NewError(fiber.StatusForbidden, err.Error())
			},
			Extractor: func(c *fiber.Ctx) (string, error) {
				if c.Get(fiber.HeaderContentType) == fiber.MIMEApplicationForm ||
					c.Get(fiber.HeaderContentType) == fiber.MIMEMultipartForm {
					return formExtractor(c)
				}
				return headerExtractor(c)
			},
			Session:           store,
			ContextKey:        "csrf",
			SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
		}
	})
}
