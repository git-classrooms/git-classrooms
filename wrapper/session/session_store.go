package session

import (
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/storage/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/memory"
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

func InitSessionStore(dsn *string, publicURL *url.URL) {
	once.Do(func() {
		var storage fiber.Storage
		if dsn == nil {
			storage = memory.New()
		} else {
			storage = postgres.New(postgres.Config{
				SslMode:       "disable",
				ConnectionURI: *dsn,
				Table:         "fiber_storage",
				Reset:         false,
				GCInterval:    10 * time.Second,
			})
		}

		var sesstionCookieName string
		var csrfCookieName string
		if publicURL.Scheme == "https" {
			sesstionCookieName = "__Host-session_id"
			csrfCookieName = "__Host-csrf_"
		} else {
			sesstionCookieName = "session_id"
			csrfCookieName = "csrf_"
		}

		store = session.New(session.Config{
			Expiration:     time.Hour * 24 * 7,
			KeyLookup:      "cookie:" + sesstionCookieName,
			CookieHTTPOnly: true,
			CookieSecure:   publicURL.Scheme == "https",
			Storage:        storage,
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
			CookieName:        csrfCookieName,
			CookieSameSite:    "Lax",
			CookieSecure:      publicURL.Scheme == "https",
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			Expiration:        1 * time.Hour,
			KeyGenerator:      utils.UUIDv4,
			// TODO: Discuss if we reactivate this in the future
			// SingleUseToken:    true,
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
