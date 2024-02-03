package router

import (
	"fmt"
	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
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
	app.Get("/auth", authController.Auth)
	app.Get(config.GetRedirectUrl().Path, authController.Callback)

	api := app.Group("/api", logger.New()) // behind "/api" is always a user logged into the session and this user is logged into the repository, which is accessable via "ctx.Locals("gitlab-repo").(repository.Repository)"
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	api.Use(authController.AuthMiddleware)

	me := api.Group("/me")
	me.Get("/", apiController.GetMe)

	me.Get("/classrooms", apiController.GetMeClassrooms)                                                                         // -> {ownClassrooms: [Classroom], joinedClassrooms: [Classroom with role]}
	me.Use("/classrooms/:classroomId", apiController.GetMeClassroomMiddleware)                                                   // puts {Classroom with role} into context
	me.Get("/classrooms/:classroomId", apiController.GetMeClassroom)                                                             // -> {Classroom with role}
	me.Get("/classrooms/:classroomId/templateProjects", apiController.GetMeClassroomTemplates)                                   // -> [Gitlab Projects public von dir oder private from classroom]
	me.Get("/classrooms/:classroomId/invitations", apiController.GetMeClassroomInvitations)                                      // -> [Invitations only owner]
	me.Get("/classrooms/:classroomId/members", apiController.GetMeClassroomMembers)                                              // -> [Members with role]
	me.Get("/classrooms/:classroomId/members/:memberId", apiController.GetMeClassroomMember)                                     // -> {Member with role}
	me.Get("/classrooms/:classroomId/members/:memberId/assignments", apiController.GetMeClassroomMemberAssignments)              // -> [assignments of member only owner Assignment with status if accepted and link to gitlab]
	me.Get("/classrooms/:classroomId/members/:memberId/assignments/:assignmentId", apiController.GetMeClassroomMemberAssignment) // -> {assignment of member only owner Assignment with status if accepted and link to gitlab}
	me.Get("/classrooms/:classroomId/assignments", apiController.GetMeClassroomAssignments)                                      // -> [Assignment with status if accepted]
	me.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetMeClassroomAssignment)                         // -> {Assignment with status if accepted and link to gitlab}

	api.Post("/classrooms", apiController.CreateClassroom)
	api.Get("/classrooms/:classroomId/assignments", apiController.GetClassroomAssignments)
	api.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetClassroomAssignment)
	api.Post("/classrooms/:classroomId/assignments", apiController.CreateAssignment)
	api.Post("/classrooms/:classroomId/members", apiController.InviteToClassroom)
	api.Post("/classrooms/:classroomId/invitations/:invitationId", apiController.JoinClassroom)
	api.Post("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignment)
	api.Post("/classrooms/:classroomId/assignments/:assignmentId/accept", apiController.JoinAssignment)

	app.Static("/", frontendPath)
	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := fmt.Sprintf("%s/index.html", frontendPath)
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
