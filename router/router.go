package router

import (
	"path"
	"strings"

	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
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
	v2Controller apiV2.Controller,
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
	setupV2Routes(&api, config, authController, v2Controller)

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

	v1.Get("/me", apiController.GetMe)
	v1.Get("/me/gitlab", apiController.GetMeGitlab)

	v1.Get("/classrooms/owned", apiController.GetOwnedClassrooms)
	v1.Post("/classrooms/owned", apiController.CreateClassroom)
	v1.Use("/classrooms/owned/:classroomId", apiController.OwnedClassroomMiddleware)
	v1.Get("/classrooms/owned/:classroomId", apiController.GetOwnedClassroom)
	v1.Put("/classrooms/owned/:classroomId", apiController.PutOwnedClassroom)
	v1.Get("/classrooms/owned/:classroomId/gitlab", apiController.RedirectGroupGitlab)

	v1.Get("/classrooms/owned/:classroomId/assignments", apiController.GetOwnedClassroomAssignments)
	v1.Post("/classrooms/owned/:classroomId/assignments", apiController.CreateAssignment)
	v1.Use("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.OwnedClassroomAssignmentMiddleware)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId", apiController.GetOwnedClassroomAssignment)

	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.GetOwnedClassroomAssignmentProjects)
	v1.Post("/classrooms/owned/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignmentProject)
	v1.Use("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.OwnedClassroomAssignmentProjectMiddleware)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.GetOwnedClassroomAssignmentProject)
	v1.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)

	v1.Get("/classrooms/owned/:classroomId/members", apiController.GetOwnedClassroomMembers)
	v1.Use("/classrooms/owned/:classroomId/members/:memberId", apiController.OwnedClassroomMemberMiddleware)
	v1.Get("/classrooms/owned/:classroomId/members/:memberId", apiController.GetOwnedClassroomMember)
	v1.Patch("/classrooms/owned/:classroomId/members/:memberId", apiController.ChangeOwnedClassroomMember)
	v1.Get("/classrooms/owned/:classroomId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	v1.Get("/classrooms/owned/:classroomId/invitations", apiController.GetOwnedClassroomInvitations)
	v1.Post("/classrooms/owned/:classroomId/invitations", apiController.InviteToClassroom)

	v1.Get("/classrooms/owned/:classroomId/templateProjects", apiController.GetOwnedClassroomTemplates)

	v1.Get("/classrooms/owned/:classroomId/teams", apiController.GetOwnedClassroomTeams)
	v1.Post("classrooms/owned/:classroomId/teams", apiController.CreateOwnedClassroomTeam)
	v1.Use("/classrooms/owned/:classroomId/teams/:teamId", apiController.OwnedClassroomTeamMiddleware)
	v1.Get("/classrooms/owned/:classroomId/teams/:teamId", apiController.GetOwnedClassroomTeam)
	v1.Get("/classrooms/owned/:classroomId/teams/:teamId/gitlab", apiController.RedirectGroupGitlab)

	v1.Get("/classrooms/owned/:classroomId/teams/:teamId/members", apiController.GetOwnedClassroomTeamMembers)
	v1.Use("/classrooms/owned/:classroomId/teams/:teamId/members/:memberId", apiController.OwnedClassroomTeamMemberMiddleware)
	v1.Use("/classrooms/owned/:classroomId/teams/:teamId/members/:memberId/gitlab", apiController.RedirectUserGitlab)
	v1.Delete("/classrooms/owned/:classroomId/teams/:teamId/members/:memberId", apiController.RemoveMemberFromTeam)

	v1.Get("/classrooms/owned/:classroomId/teams/:teamId/projects", apiController.GetOwnedClassroomTeamProjects)

	v1.Get("/invitations/:invitationId", apiController.GetInvitationInfo)

	v1.Get("/classrooms/joined", apiController.GetJoinedClassrooms)
	v1.Post("/classrooms/joined", apiController.JoinClassroom) // with invitation id in the body
	v1.Use("/classrooms/joined/:classroomId", apiController.JoinedClassroomMiddleware)
	v1.Get("/classrooms/joined/:classroomId", apiController.GetJoinedClassroom)
	v1.Get("/classrooms/joined/:classroomId/gitlab", apiController.RedirectGroupGitlab)

	v1.Get("/classrooms/joined/:classroomId/assignments", apiController.GetJoinedClassroomAssignments)
	v1.Use("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.JoinedClassroomAssignmentMiddleware)
	v1.Get("/classrooms/joined/:classroomId/assignments/:assignmentId", apiController.GetJoinedClassroomAssignment)
	v1.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/accept", apiController.AcceptAssignment)

	v1.Get("/classrooms/joined/:classroomId/teams", apiController.GetJoinedClassroomTeams)
	v1.Post("/classrooms/joined/:classroomId/teams", apiController.CreateJoinedClassroomTeam)
	v1.Use("/classrooms/joined/:classroomId/teams/:teamId", apiController.JoinedClassroomTeamMiddleware)
	v1.Get("/classrooms/joined/:classroomId/teams/:teamId", apiController.GetJoinedClassroomTeam)
	v1.Get("/classrooms/joined/:classroomId/teams/:teamId/giltab", apiController.RedirectGroupGitlab)
	v1.Post("/classrooms/joined/:classroomId/teams/:teamId/join", apiController.JoinJoinedClassroomTeam)

	// api.Post("/classrooms/joined/:classroomId/assignments/:assignmentId/projects", apiController.InviteToAssignment) // moderator only
	// api.Post("/classrooms/joined/:classroomId/invitations", apiController.InviteToClassroom) // moderator only
	//

}

func setupV2Routes(api *fiber.Router, config authConfig.Config, authController authController.Controller,
	apiController apiV2.Controller) {

	v2 := (*api).Group("/v2")

	v2.Get("/info/gitlab", apiController.GetGitlabInfo)
	v2.Post("/auth/sign-in", authController.SignIn)
	v2.Post("/auth/sign-out", authController.SignOut)
	v2.Get(strings.Replace(config.GetRedirectUrl().Path, "/api/v2", "", 1), authController.Callback)
	v2.Get("/auth/csrf", authController.GetCsrf)
	v2.Use(authController.AuthMiddleware)
	v2.Get("/auth", authController.GetAuth)

	v2.Get("/me", apiController.GetMe)
	v2.Get("/me/gitlab", apiController.GetMeGitlab)

	v2.Get("/classrooms", apiController.GetClassrooms)
	v2.Post("/classrooms", apiController.CreateClassroom)

	v2.Get("/classrooms/:classroomId/invitations/:invitationId", apiController.GetClassroomInvitation)
	v2.Post("/classrooms/:classroomId/join", apiController.JoinClassroom) // with invitation id in the body

	v2.Use("/classrooms/:classroomId", apiController.ClassroomMiddleware, apiController.ArchivedMiddleware)
	v2.Get("/classrooms/:classroomId", apiController.GetClassroom)
	v2.Put("/classrooms/:classroomId", apiController.CreatorMiddleware(), apiController.UpdateClassroom)
	v2.Patch("/classrooms/:classroomId/archive", apiController.CreatorMiddleware(), apiController.ArchiveClassroom)
	v2.Get("/classrooms/:classroomId/gitlab", apiController.RedirectGroupGitlab)

	v2.Get("/classrooms/:classroomId/templateProjects", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomTemplates)

	v2.Use("/classrooms/:classroomId/assignments", apiController.ViewableClassroomMiddleware())
	v2.Get("/classrooms/:classroomId/assignments", apiController.GetClassroomAssignments)
	v2.Post("/classrooms/:classroomId/assignments", apiController.RoleMiddleware(database.Owner), apiController.CreateAssignment)
	v2.Use("/classrooms/:classroomId/assignments/:assignmentId", apiController.ClassroomAssignmentMiddleware)
	v2.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetClassroomAssignment)
	v2.Put("/classrooms/:classroomId/assignments/:assignmentId", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignment)

	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.GetGradingRubrics)
	v2.Put("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.UpdateGradingRubrics)
	v2.Post("/classrooms/:classroomId/assignments/:assignmentId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGrading)

	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/repos", apiController.GetMultipleProjectCloneUrls)
	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.GetClassroomAssignmentProjects)
	v2.Post("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.RoleMiddleware(database.Owner), apiController.InviteToAssignment)
	v2.Use("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.ClassroomAssignmentProjectMiddleware)
	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.GetClassroomAssignmentProject)

	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.GetGradingResults)
	v2.Put("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.UpdateGradingResults)
	v2.Post("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGradingForProject)

	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	v2.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/repo", apiController.GetProjectCloneUrls)

	v2.Use("/classrooms/:classroomId/projects", apiController.RoleMiddleware(database.Student))
	v2.Get("/classrooms/:classroomId/projects", apiController.GetClassroomProjects)
	v2.Use("/classrooms/:classroomId/projects/:projectId", apiController.ClassroomProjectMiddleware)
	v2.Get("/classrooms/:classroomId/projects/:projectId", apiController.GetClassroomProject)
	v2.Post("/classrooms/:classroomId/projects/:projectId/accept", apiController.AcceptAssignment)
	v2.Get("/classrooms/:classroomId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	v2.Get("/classrooms/:classroomId/projects/:projectId/repo", apiController.GetProjectCloneUrls)
	v2.Get("/classrooms/:classroomId/projects/:projectId/grading", apiController.GetGradingResults)

	v2.Use("/classrooms/:classroomId/invitations", apiController.RoleMiddleware(database.Owner, database.Moderator))
	v2.Get("/classrooms/:classroomId/invitations", apiController.GetClassroomInvitations)
	v2.Post("/classrooms/:classroomId/invitations", apiController.InviteToClassroom)
	v2.Delete("/classrooms/:classroomId/invitations/:invitationId", apiController.RevokeClassroomInvitation)

	v2.Get("/classrooms/:classroomId/members", apiController.GetClassroomMembers)
	v2.Use("/classrooms/:classroomId/members/:memberId", apiController.ClassroomMemberMiddleware)
	v2.Get("/classrooms/:classroomId/members/:memberId", apiController.GetClassroomMember)
	v2.Patch("/classrooms/:classroomId/members/:memberId/team", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.UpdateMemberTeam)
	v2.Patch("/classrooms/:classroomId/members/:memberId/role", apiController.RoleMiddleware(database.Owner), apiController.UpdateMemberRole)
	// v2.Delete("/classrooms/:classroomId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveClassroomMember)
	v2.Get("/classrooms/:classroomId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	v2.Get("/classrooms/:classroomId/runners", apiController.GetClassroomRunners)
	v2.Get("/classrooms/:classroomId/runners/available", apiController.GetClassroomRunnersAreAvailable)

	v2.Get("/classrooms/:classroomId/teams", apiController.GetClassroomTeams)
	v2.Post("/classrooms/:classroomId/teams", apiController.CreateTeam)
	v2.Use("/classrooms/:classroomId/teams/:teamId", apiController.ClassroomTeamMiddleware)
	v2.Get("/classrooms/:classroomId/teams/:teamId", apiController.GetClassroomTeam)
	v2.Put("/classrooms/:classroomId/teams/:teamId", apiController.UpdateTeam)
	v2.Post("/classrooms/:classroomId/teams/:teamId/join", apiController.RoleMiddleware(database.Student), apiController.JoinTeam)
	v2.Get("/classrooms/:classroomId/teams/:teamId/gitlab", apiController.RedirectGroupGitlab)

	v2.Get("/classrooms/:classroomId/teams/:teamId/members", apiController.GetClassroomTeamMembers)
	v2.Use("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.ClassroomTeamMemberMiddleware)
	v2.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.GetClassroomTeamMember)
	v2.Delete("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveMemberFromTeam)
	v2.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	v2.Use("/classrooms/:classroomId/teams/:teamId/projects", apiController.ViewableClassroomMiddleware())
	v2.Get("/classrooms/:classroomId/teams/:teamId/projects", apiController.GetClassroomTeamProjects)
	v2.Use("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.ClassroomTeamProjectMiddleware)
	v2.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.GetClassroomTeamProject)
	v2.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
}

func setupFrontend(app *fiber.App, frontendPath string) {
	app.Static("/", frontendPath)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := path.Join(frontendPath, "index.html")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
