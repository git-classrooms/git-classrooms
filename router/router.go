package router

import (
	"fmt"
	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(
	app *fiber.App,
	authController authController.Controller,
	apiController apiController.Controller,
	frontendPath string,
	config authConfig.Config,
) {
	app.Static("/", frontendPath)

	app.Get("/auth", authController.Auth)
	app.Get(config.GetRedirectUrl().Path, authController.Callback)

	api := app.Group("/api", logger.New()) // behind "/api" is always a user logged into the session and this user is logged into the repository, which is accessable via "ctx.Locals("gitlab-repo").(repository.Repository)"
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	api.Use(authController.AuthMiddleware)

	api.Get("/secret", func(c *fiber.Ctx) error {
		repo := context.GetGitlabRepository(c)
		user, err := repo.GetCurrentUser()
		if err != nil {
			return err
		}

		return c.JSON(user)
	})

	api.Post("/classrooms", apiController.CreateClassroom)
	api.Post("/assignments", apiController.CreateAssignment)
	api.Post("/classrooms/:classroomId/members", apiController.JoinClassroom)
	api.Post("/classrooms/:classroomId/invitations/:inviteId", apiController.InviteToClassroom)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := fmt.Sprintf("%s/index.html", frontendPath)
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
