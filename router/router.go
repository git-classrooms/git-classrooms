package router

import (
	"path"

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
	app.Get("/auth", authController.SignIn)
	app.Post("/auth/logout", authController.SignOut)
	api := app.Group("/api", logger.New()) // behind "/api" is always a user logged into the session and this user is logged into the repository, which is accessable via "ctx.Locals("gitlab-repo").(repository.Repository)"

	api.Get("/auth/sign-in", authController.SignIn)
	api.Post("/auth/sign-out", authController.SignOut)
	app.Get(config.GetRedirectUrl().Path, authController.Callback)
	api.Use(authController.AuthMiddleware)
	api.Get("/auth", authController.GetAuth)

	//	"/me"
	//	"/users/:userId/classrooms/" -> middleware
	//
	//	"/classrooms/joined/:classroomId"
	//
	//	"/ownedClassrooms/members/assignment"
	//	"/joinedClassrooms/members/assignments"

	api.Get("/me", apiController.GetMe)

	api.Get("/classrooms/owned", apiController.GetOwnedClassrooms)
	api.Post("/classrooms/owned", apiController.CreateClassroom)
	api.Use("/classrooms/owned/:classroomId", apiController.OwnedClassroomMiddleware)
	api.Get("/classrooms/owned/:classroomId", apiController.GetOwnedClassroom)

	api.Get("/classrooms/owned/:classroomId/assignments", apiController.GetOwnedClassroomAssignments)
	api.Post("/classrooms/owned/:classroomId/assignments", apiController.CreateAssignment)
	api.Use("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.OwnedClassroomAssignmentMiddleware)
	api.Get("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.GetOwnedClassroomAssignment)

	api.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.GetClassroomAssignmentProjects)
	api.Post("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignmentProject)

	api.Get("/classrooms/owned/:classroomId/members", apiController.GetOwnedClassroomMembers)

	api.Get("/classrooms/owned/:classroomId/invitations", apiController.GetOwnedClassroomInvitations)
	api.Post("/classrooms/owned/:classroomId/invitations", apiController.InviteToClassroom)

	api.Get("/classrooms/owned/:classroomId/templateProjects", apiController.GetOwnedClassroomTemplates)

	api.Get("/classrooms/joined", apiController.GetJoinedClassrooms)
	api.Post("/classrooms/joined", apiController.JoinClassroomNew) // with invitation id in the body
	api.Use("/classrooms/joined/:classroomId", apiController.JoinedClassroomMiddleware)
	api.Get("/classrooms/joined/:classroomId", apiController.GetJoinedClassroom)

	api.Get("/classrooms/joined/:classroomId/assignments", apiController.GetJoinedClassroomAssignments)
	api.Use("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.JoinedClassroomAssignmentMiddleware)
	api.Get("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.GetJoinedClassroomAssignment)
	api.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/accept", apiController.JoinAssignmentNew)

	// api.Get("/classrooms/owned/:classroomId/members/:memberId", apiController.GetOwnedClassroomMember)
	// api.Get("/classrooms/owned/:classroomId/members/:memberId/assignments", apiController.GetOwnedClassroomMemberAssignments)
	// api.Get("/classrooms/owned/:classroomId/members/:memberId/assignments/:assignmentId", apiController.GetOwnedClassroomMemberAssignment)
	//
	// api.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignment) // moderator only
	// api.Post("/classrooms/joined/:classroomId/invitations", apiController.InviteToClassroom) // moderator only
	//

	me := api.Group("/me")
	me.Get("/", apiController.GetMe)

	// TODO: Nochmal überlegen ob das alles so gut is wie wir es jetzt tun
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

	// TODO: Namen angleichen der Dateien
	// TODO: Get current invitations from classroom
	api.Post("/classrooms", apiController.CreateClassroom)
	api.Get("/classrooms/:classroomId/assignments", apiController.GetClassroomAssignments)
	api.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetClassroomAssignment)
	api.Post("/classrooms/:classroomId/assignments", apiController.CreateAssignment)
	api.Post("/classrooms/:classroomId/members", apiController.InviteToClassroom)
	api.Post("/classrooms/:classroomId/invitations/:invitationId", apiController.JoinClassroom)
	api.Post("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignment)
	api.Get("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.GetClassroomAssignmentProjects)
	api.Post("/classrooms/:classroomId/assignments/:assignmentId/accept", apiController.JoinAssignment)

	setupFrontend(app, frontendPath)
}

func setupFrontend(app *fiber.App, frontendPath string) {
	app.Static("/", frontendPath)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := path.Join(frontendPath, "index.html")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
