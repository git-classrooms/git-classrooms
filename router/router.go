package router

import (
	"io/fs"
	"net/http"
	"strings"

	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "gitlab.hs-flensburg.de/gitlab-classroom/docs"
)

func Routes(
	authController authController.Controller,
	apiController apiController.Controller,
	frontendFS fs.FS,
	config authConfig.Config,
) *fiber.App {
	app := fiber.New()
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

	api := app.Group("/api", logger.New())
	api.Mount("/v1", setupApiRoutes(config, authController, apiController))
	api.Get("/swagger/*", swagger.HandlerDefault) // default
	api.Get("/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	app.Use("/", filesystem.New(filesystem.Config{
		PathPrefix:   "frontend/dist",
		Index:        "index.html",
		NotFoundFile: "frontend/dist/index.html",
		Root:         http.FS(frontendFS),
	}))

	return app
}

func setupApiRoutes(config authConfig.Config, authController authController.Controller,
	apiController apiController.Controller) *fiber.App {

	app := fiber.New()

	app.Get("/info/gitlab", apiController.GetGitlabInfo)
	app.Post("/auth/sign-in", authController.SignIn)
	app.Post("/auth/sign-out", authController.SignOut)
	app.Get(strings.Replace(config.GetRedirectUrl().Path, "/api/v1", "", 1), authController.Callback)
	app.Get("/auth/csrf", authController.GetCsrf)
	app.Use(authController.AuthMiddleware)
	app.Get("/auth", authController.GetAuth)

	app.Get("/me", apiController.GetMe)
	app.Get("/me/gitlab", apiController.GetMeGitlab)

	app.Get("/assignments", apiController.GetActiveAssignments)

	app.Get("/classrooms", apiController.GetClassrooms)
	app.Post("/classrooms", apiController.CreateClassroom)

	app.Get("/classrooms/:classroomId/invitations/:invitationId", apiController.GetClassroomInvitation)
	app.Post("/classrooms/:classroomId/join", apiController.JoinClassroom) // with invitation id in the body

	app.Use("/classrooms/:classroomId", apiController.ClassroomMiddleware, apiController.PotentiallyDeletedClassroomMiddleware, apiController.ArchivedMiddleware, apiController.RotateAccessTokenMiddleware)
	app.Get("/classrooms/:classroomId", apiController.GetClassroom)
	app.Put("/classrooms/:classroomId", apiController.CreatorMiddleware(), apiController.UpdateClassroom)
	app.Patch("/classrooms/:classroomId/archive", apiController.CreatorMiddleware(), apiController.ArchiveClassroom)
	app.Get("/classrooms/:classroomId/gitlab", apiController.RedirectGroupGitlab)

	app.Get("/classrooms/:classroomId/grading", apiController.RoleMiddleware(database.Owner), apiController.GetGradingRubrics)
	app.Put("/classrooms/:classroomId/grading", apiController.RoleMiddleware(database.Owner), apiController.UpdateGradingRubrics)
	app.Get("/classrooms/:classroomId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomReport)

	app.Get("/classrooms/:classroomId/templateProjects", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomTemplates)

	app.Use("/classrooms/:classroomId/assignments", apiController.ViewableClassroomMiddleware())
	app.Get("/classrooms/:classroomId/assignments", apiController.GetClassroomAssignments)
	app.Post("/classrooms/:classroomId/assignments", apiController.RoleMiddleware(database.Owner), apiController.CreateAssignment)
	app.Use("/classrooms/:classroomId/assignments/:assignmentId", apiController.ClassroomAssignmentMiddleware)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId", apiController.GetClassroomAssignment)
	app.Put("/classrooms/:classroomId/assignments/:assignmentId", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignment)

	app.Get("/classrooms/:classroomId/assignments/:assignmentId/tests", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomAssignmentTests)
	app.Put("/classrooms/:classroomId/assignments/:assignmentId/tests", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignmentTests)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.GetAssignmentGradingRubrics)
	app.Put("/classrooms/:classroomId/assignments/:assignmentId/grading", apiController.RoleMiddleware(database.Owner), apiController.UpdateAssignmentGradingRubrics)
	app.Post("/classrooms/:classroomId/assignments/:assignmentId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGrading)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomAssignmentReport)

	app.Get("/classrooms/:classroomId/assignments/:assignmentId/repos", apiController.GetMultipleProjectCloneUrls)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.GetClassroomAssignmentProjects)
	app.Post("/classrooms/:classroomId/assignments/:assignmentId/projects", apiController.RoleMiddleware(database.Owner), apiController.InviteToAssignment)
	app.Use("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.ClassroomAssignmentProjectMiddleware)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", apiController.GetClassroomAssignmentProject)

	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.GetGradingResults)
	app.Put("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.UpdateGradingResults)
	app.Post("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/grading/auto", apiController.RoleMiddleware(database.Owner, database.Moderator), apiController.StartAutoGradingForProject)

	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/report/gitlab", apiController.RedirectReportGitlab)
	app.Get("/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId/repo", apiController.GetProjectCloneUrls)

	app.Use("/classrooms/:classroomId/projects", apiController.RoleMiddleware(database.Student))
	app.Get("/classrooms/:classroomId/projects", apiController.GetClassroomProjects)
	app.Use("/classrooms/:classroomId/projects/:projectId", apiController.ClassroomProjectMiddleware)
	app.Get("/classrooms/:classroomId/projects/:projectId", apiController.GetClassroomProject)
	app.Post("/classrooms/:classroomId/projects/:projectId/accept", apiController.AcceptAssignment)
	app.Get("/classrooms/:classroomId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	app.Get("/classrooms/:classroomId/projects/:projectId/report/gitlab", apiController.RedirectProjectGitlab)
	app.Get("/classrooms/:classroomId/projects/:projectId/repo", apiController.GetProjectCloneUrls)
	app.Get("/classrooms/:classroomId/projects/:projectId/grading", apiController.GetGradingResults)

	app.Use("/classrooms/:classroomId/invitations", apiController.RoleMiddleware(database.Owner, database.Moderator))
	app.Get("/classrooms/:classroomId/invitations", apiController.GetClassroomInvitations)
	app.Post("/classrooms/:classroomId/invitations", apiController.InviteToClassroom)
	app.Delete("/classrooms/:classroomId/invitations/:invitationId", apiController.RevokeClassroomInvitation)

	app.Get("/classrooms/:classroomId/members", apiController.GetClassroomMembers)
	app.Use("/classrooms/:classroomId/members/:memberId", apiController.ClassroomMemberMiddleware)
	app.Get("/classrooms/:classroomId/members/:memberId", apiController.GetClassroomMember)
	app.Patch("/classrooms/:classroomId/members/:memberId/team", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.UpdateMemberTeam)
	app.Patch("/classrooms/:classroomId/members/:memberId/role", apiController.RoleMiddleware(database.Owner), apiController.UpdateMemberRole)
	// app.Delete("/classrooms/:classroomId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveClassroomMember)
	app.Get("/classrooms/:classroomId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	app.Get("/classrooms/:classroomId/runners", apiController.GetClassroomRunners)
	app.Get("/classrooms/:classroomId/runners/available", apiController.GetClassroomRunnersAreAvailable)

	app.Get("/classrooms/:classroomId/teams", apiController.GetClassroomTeams)
	app.Post("/classrooms/:classroomId/teams", apiController.CreateTeam)
	app.Use("/classrooms/:classroomId/teams/:teamId", apiController.ClassroomTeamMiddleware)
	app.Get("/classrooms/:classroomId/teams/:teamId", apiController.GetClassroomTeam)
	app.Put("/classrooms/:classroomId/teams/:teamId", apiController.UpdateTeam)
	app.Post("/classrooms/:classroomId/teams/:teamId/join", apiController.RoleMiddleware(database.Student), apiController.JoinTeam)
	app.Get("/classrooms/:classroomId/teams/:teamId/gitlab", apiController.RedirectGroupGitlab)

	app.Get("/classrooms/:classroomId/teams/:teamId/members", apiController.GetClassroomTeamMembers)
	app.Use("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.ClassroomTeamMemberMiddleware)
	app.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.GetClassroomTeamMember)
	app.Delete("/classrooms/:classroomId/teams/:teamId/members/:memberId", apiController.RoleMiddleware(database.Moderator, database.Owner), apiController.RemoveMemberFromTeam)
	app.Get("/classrooms/:classroomId/teams/:teamId/members/:memberId/gitlab", apiController.RedirectUserGitlab)

	app.Use("/classrooms/:classroomId/teams/:teamId/projects", apiController.ViewableClassroomMiddleware())
	app.Get("/classrooms/:classroomId/teams/:teamId/projects", apiController.GetClassroomTeamProjects)
	app.Use("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.ClassroomTeamProjectMiddleware)
	app.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId", apiController.GetClassroomTeamProject)
	app.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId/gitlab", apiController.RedirectProjectGitlab)
	app.Get("/classrooms/:classroomId/teams/:teamId/projects/:projectId/report/gitlab", apiController.RedirectReportGitlab)

	app.Get("/classrooms/:classroomId/teams/:teamId/grading/report", apiController.RoleMiddleware(database.Owner), apiController.GetClassroomTeamReport)

	return app
}
