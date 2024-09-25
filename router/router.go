package router

import (
	"path"
	"strings"

	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
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
	setupRoutes(&api, config, authController, apiController)

	api.Get("/swagger/*", swagger.HandlerDefault) // default

	setupFrontend(app, frontendPath)
}

func setupRoutes(api *fiber.Router, config authConfig.Config, authController authController.Controller,
	apiController apiController.Controller) {

	v1 := (*api).Group("/v1")

	v1.Get("/info/gitlab", apiController.GetGitlabInfo)
	v1.Post("/auth/sign-in", authController.SignIn)
	v1.Post("/auth/sign-out", authController.SignOut)
	v1.Get(strings.Replace(config.GetRedirectUrl().Path, "/api/v1", "", 1), authController.Callback)
	v1.Get("/auth/csrf", authController.GetCsrf)
	v1.Use(authController.AuthMiddleware)
	v1.Get("/auth", authController.GetAuth)

	v1.Get("/me", apiController.GetMe)
	v1.Get("/me/gitlab", apiController.GetMeGitlab)

	v1.Get("/assignments", apiController.GetActiveAssignments)

	v1.Get("/classrooms", apiController.GetClassrooms)
	v1.Post("/classrooms", apiController.CreateClassroom)

	v1.Get("/classrooms/:classroomId/invitations/:invitationId", apiController.GetClassroomInvitation)
	v1.Post("/classrooms/:classroomId/join", apiController.JoinClassroom) // with invitation id in the body

	v1.Use("/classrooms/:classroomId", apiController.ClassroomMiddleware, apiController.PotentiallyDeletedClassroomMiddleware, apiController.ArchivedMiddleware, apiController.RotateAccessTokenMiddleware)
	v1.Get("/classrooms/:classroomId", apiController.GetClassroom)
	v1.Put("/classrooms/:classroomId", apiController.CreatorMiddleware(), apiController.UpdateClassroom)
	v1.Patch("/classrooms/:classroomId/archive", apiController.CreatorMiddleware(), apiController.ArchiveClassroom)
	v1.Get("/classrooms/:classroomId/gitlab", apiController.RedirectGroupGitlab)

	v1.Get("/classrooms/:classroomId/grading", apiController.RoleMiddleware(database.Owner), apiController.GetGradingRubrics)
	v1.Put("/classrooms/:classroomId/grading", apiController.RoleMiddleware(database.Owner), apiController.UpdateGradingRubrics)
	v1.Get("/classrooms/:classroomId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomReport)

	v1.Get("/classrooms/:classroomId/templateProjects", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomTemplates)

	v1.Use("/classrooms/:classroomId/assignments", apiController.ViewableClassroomMiddleware())
	v1.Get("/classrooms/:classroomId/assignments", apiController.GetClassroomAssignments)
	v1.Post("/classrooms/:classroomId/assignments", apiController.RoleMiddleware(database.Owner), apiController.CreateAssignment)
	v1.Use("/classrooms/:classroomId/assignments/:assignmentId", apiController.ClassroomAssignmentMiddleware)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetClassroomAssignment)
	v1.Put("/classrooms/:classroomId/assignments/:assignmentId", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignment)

	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/tests", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomAssignmentTests)
	v1.Put("/classrooms/:classroomId/assignments/:assignmentId/tests", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignmentTests)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.GetAssignmentGradingRubrics)
	v1.Put("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignmentGradingRubrics)
	v1.Post("/classrooms/:classroomId/assignments/:assignmentId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGrading)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomAssignmentReport)

	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/repos", apiController.GetMultipleProjectCloneUrls)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.GetClassroomAssignmentProjects)
	v1.Post("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.RoleMiddleware(database.Owner), apiController.InviteToAssignment)
	v1.Use("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.ClassroomAssignmentProjectMiddleware)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.GetClassroomAssignmentProject)

	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.GetGradingResults)
	v1.Put("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.UpdateGradingResults)
	v1.Post("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGradingForProject)

	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/report/gitlab", apiController.RedirectReportGitlab)
	v1.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/repo", apiController.GetProjectCloneUrls)

	v1.Use("/classrooms/:classroomId/projects", apiController.RoleMiddleware(database.Student))
	v1.Get("/classrooms/:classroomId/projects", apiController.GetClassroomProjects)
	v1.Use("/classrooms/:classroomId/projects/:projectId", apiController.ClassroomProjectMiddleware)
	v1.Get("/classrooms/:classroomId/projects/:projectId", apiController.GetClassroomProject)
	v1.Post("/classrooms/:classroomId/projects/:projectId/accept", apiController.AcceptAssignment)
	v1.Get("/classrooms/:classroomId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	v1.Get("/classrooms/:classroomId/projects/:projectId/report/gitlab", apiController.RedirectProjectGitlab)
	v1.Get("/classrooms/:classroomId/projects/:projectId/repo", apiController.GetProjectCloneUrls)
	v1.Get("/classrooms/:classroomId/projects/:projectId/grading", apiController.GetGradingResults)

	v1.Use("/classrooms/:classroomId/invitations", apiController.RoleMiddleware(database.Owner, database.Moderator))
	v1.Get("/classrooms/:classroomId/invitations", apiController.GetClassroomInvitations)
	v1.Post("/classrooms/:classroomId/invitations", apiController.InviteToClassroom)
	v1.Delete("/classrooms/:classroomId/invitations/:invitationId", apiController.RevokeClassroomInvitation)

	v1.Get("/classrooms/:classroomId/members", apiController.GetClassroomMembers)
	v1.Use("/classrooms/:classroomId/members/:memberId", apiController.ClassroomMemberMiddleware)
	v1.Get("/classrooms/:classroomId/members/:memberId", apiController.GetClassroomMember)
	v1.Patch("/classrooms/:classroomId/members/:memberId/team", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.UpdateMemberTeam)
	v1.Patch("/classrooms/:classroomId/members/:memberId/role", apiController.RoleMiddleware(database.Owner), apiController.UpdateMemberRole)
	// v1.Delete("/classrooms/:classroomId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveClassroomMember)
	v1.Get("/classrooms/:classroomId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	v1.Get("/classrooms/:classroomId/runners", apiController.GetClassroomRunners)
	v1.Get("/classrooms/:classroomId/runners/available", apiController.GetClassroomRunnersAreAvailable)

	v1.Get("/classrooms/:classroomId/teams", apiController.GetClassroomTeams)
	v1.Post("/classrooms/:classroomId/teams", apiController.CreateTeam)
	v1.Use("/classrooms/:classroomId/teams/:teamId", apiController.ClassroomTeamMiddleware)
	v1.Get("/classrooms/:classroomId/teams/:teamId", apiController.GetClassroomTeam)
	v1.Put("/classrooms/:classroomId/teams/:teamId", apiController.UpdateTeam)
	v1.Post("/classrooms/:classroomId/teams/:teamId/join", apiController.RoleMiddleware(database.Student), apiController.JoinTeam)
	v1.Get("/classrooms/:classroomId/teams/:teamId/gitlab", apiController.RedirectGroupGitlab)

	v1.Get("/classrooms/:classroomId/teams/:teamId/members", apiController.GetClassroomTeamMembers)
	v1.Use("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.ClassroomTeamMemberMiddleware)
	v1.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.GetClassroomTeamMember)
	v1.Delete("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveMemberFromTeam)
	v1.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	v1.Use("/classrooms/:classroomId/teams/:teamId/projects", apiController.ViewableClassroomMiddleware())
	v1.Get("/classrooms/:classroomId/teams/:teamId/projects", apiController.GetClassroomTeamProjects)
	v1.Use("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.ClassroomTeamProjectMiddleware)
	v1.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.GetClassroomTeamProject)
	v1.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	v1.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId/report/gitlab", apiController.RedirectReportGitlab)

	v1.Get("/classrooms/:classroomId/teams/:teamId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomTeamReport)
}

func setupFrontend(app *fiber.App, frontendPath string) {
	app.Static("/", frontendPath)

	// Catch all routes
	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := path.Join(frontendPath, "index.html")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })
}
