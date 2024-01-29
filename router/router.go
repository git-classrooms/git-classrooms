package router

import (
	"backend/config"
	"backend/handler"
	apiHandler "backend/handler/api"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/auth", handler.Auth)
	app.Get(config.GetConfig().Auth.RedirectURL.Path, handler.Callback)

	api := app.Group("/api", logger.New()) // behind "/api" is always a user logged into the session and this user is logged into the repository, which is accessable via "ctx.Locals("gitlab-repo").(repository.Repository)"

	fiberHandler := apiHandler.NewFiberApiHandler()

	api.Post("/createClassroom", fiberHandler.CreateClassroom)
	api.Post("/createAssignment", fiberHandler.CreateAssignment)
}
