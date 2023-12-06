package router

import (
	"backend/config"
	"backend/handler"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/auth/", handler.Auth)
	app.Get(config.GetConfig().Auth.RedirectURL, handler.Callback)
}
