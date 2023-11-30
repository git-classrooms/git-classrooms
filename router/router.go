package router

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {
	repo := go_gitlab_repo.NewGoGitlabRepo()
	fiberHandler := handler.NewFiberHandler(repo)

	api := app.Group("/api", logger.New())

	api.Post("/createClassroom", fiberHandler.CreateClassroom)
	api.Get("/", handler.Auth)
	api.Get("/auth/gitlab/callback", handler.Callback)
}
