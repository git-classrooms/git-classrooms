package router

import (
	"path"
	"strings"

	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "gitlab.hs-flensburg.de/gitlab-classroom/docs"
)

func Routes(
	app *fiber.App,
	authController authController.Controller,
	apiController apiController.Controller,
	frontendPath string,
	config authConfig.Config,
) {
	// Init session on every request if not present
	app.Use(func(c *fiber.Ctx) error {
		sess := session.Get(c)
		if sess.Session.Fresh() {
			err := sess.Save()
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		return c.Next()
	})

	app.Use(csrf.New(session.CsrfConfig))

	api := app.Group("/api", logger.New()) // behind "/api" is always a user logged into the session and this user is logged into the repository, which is accessable via "ctx.Locals("gitlab-repo").(repository.Repository)"
	setupV1Routes(&api, config, authController, apiController)

	api.Get("/swagger/*", swagger.HandlerDefault) // default

	setupFrontend(app, frontendPath)
}

func setupV1Routes(api *fiber.Router, config authConfig.Config, authController authController.Controller,
	apiController apiController.Controller) {

	v1 := (*api).Group("/v1")

	v1.Post("/auth/sign-in", authController.SignIn)
	v1.Post("/auth/sign-out", authController.SignOut)
	v1.Get(strings.Replace(config.GetRedirectUrl().Path, "/api/v1", "", 1), authController.Callback)
	v1.Get("/auth/csrf", authController.GetCsrf)
	v1.Use(authController.AuthMiddleware)
	v1.Get("/auth", authController.GetAuth)

	//	"/me"
	//	"/users/:userId/classrooms/" -> middleware
	//
	//	"/classrooms/joined/:classroomId"
	//
	//	"/ownedClassrooms/members/assignment"
	//	"/joinedClassrooms/members/assignments"

	v1.Get("/me", apiController.GetMe)

	v1.Get("/classrooms/owned", apiController.GetOwnedClassrooms)
	v1.Post("/classrooms/owned", apiController.CreateClassroom)
	v1.Use("/classrooms/owned/:classroomId", apiController.OwnedClassroomMiddleware)
	v1.Get("/classrooms/owned/:classroomId", apiController.GetOwnedClassroom)
	v1.Get("/classrooms/owned/:classroomId/gitlab", apiController.GetOwnedClassroomGitlab)

	v1.Get("/classrooms/owned/:classroomId/assignments", apiController.GetOwnedClassroomAssignments)
	v1.Post("/classrooms/owned/:classroomId/assignments", apiController.CreateAssignment)
	v1.Use("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.OwnedClassroomAssignmentMiddleware)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.GetOwnedClassroomAssignment)

	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.GetOwnedClassroomAssignmentProjects)
	v1.Post("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignmentProject)
	v1.Use("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.OwnedClassroomAssignmentProjectMiddleware)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.GetOwnedClassroomAssignmentProject)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId/gitlab", apiController.GetOwnedClassroomAssignmentProjectGitlab)

	v1.Get("/classrooms/owned/:classroomId/members", apiController.GetOwnedClassroomMembers)

	v1.Get("/classrooms/owned/:classroomId/invitations", apiController.GetOwnedClassroomInvitations)
	v1.Post("/classrooms/owned/:classroomId/invitations", apiController.InviteToClassroom)

	v1.Get("/classrooms/owned/:classroomId/templateProjects", apiController.GetOwnedClassroomTemplates)

	v1.Get("/classrooms/joined", apiController.GetJoinedClassrooms)
	v1.Post("/classrooms/joined", apiController.JoinClassroom) // with invitation id in the body
	v1.Use("/classrooms/joined/:classroomId", apiController.JoinedClassroomMiddleware)
	v1.Get("/classrooms/joined/:classroomId", apiController.GetJoinedClassroom)
	v1.Get("/classrooms/joined/:classroomId/gitlab", apiController.GetJoinedClassroomGitlab)

	v1.Get("/classrooms/joined/:classroomId/assignments", apiController.GetJoinedClassroomAssignments)
	v1.Use("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.JoinedClassroomAssignmentMiddleware)
	v1.Get("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.GetJoinedClassroomAssignment)
	v1.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/accept", apiController.JoinAssignment)

	// api.Get("/classrooms/owned/:classroomId/members/:memberId", apiController.GetOwnedClassroomMember)
	// api.Get("/classrooms/owned/:classroomId/members/:memberId/assignments", apiController.GetOwnedClassroomMemberAssignments)
	// api.Get("/classrooms/owned/:classroomId/members/:memberId/assignments/:assignmentId", apiController.GetOwnedClassroomMemberAssignment)
	//
	// api.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignment) // moderator only
	// api.Post("/classrooms/joined/:classroomId/invitations", apiController.InviteToClassroom) // moderator only
	//

}

func setupFrontend(app *fiber.App, frontendPath string) {
	app.Static("/", frontendPath)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := path.Join(frontendPath, "index.html")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
