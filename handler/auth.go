package handler

import (
	"backend/auth"

	"github.com/gofiber/fiber/v2"
)

// Auth fiber handler
func Auth(c *fiber.Ctx) error {
	path := auth.ConfigGitlab()
	url := path.AuthCodeURL("state")
	return c.Redirect(url)

}

// Callback to receive gitlabs's response
func Callback(c *fiber.Ctx) error {
	token, error := auth.ConfigGitlab().Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	return c.Status(200).JSON(fiber.Map{"token": token})

}
