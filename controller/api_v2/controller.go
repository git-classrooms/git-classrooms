package apiV2

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

type ValidateUserFunc func(user database.UserClassrooms) bool

type Controller interface {
	ValidateUserMiddleware(ValidateUserFunc) fiber.Handler
	RoleMiddleware(validRoles ...database.Role) fiber.Handler
	CreatorMiddleware() fiber.Handler
	ViewableClassroomMiddleware() fiber.Handler

	ArchivedMiddleware(*fiber.Ctx) error

	RedirectUserGitlab(*fiber.Ctx) error
	RedirectGroupGitlab(*fiber.Ctx) error
	RedirectProjectGitlab(*fiber.Ctx) error

	GetProjectCloneUrls(*fiber.Ctx) error
	GetMultipleProjectCloneUrls(*fiber.Ctx) error

	GetMe(*fiber.Ctx) error
	GetMeGitlab(*fiber.Ctx) error
	GetActiveAssignments(*fiber.Ctx) error

	GetClassrooms(*fiber.Ctx) error
	CreateClassroom(*fiber.Ctx) error
	ClassroomMiddleware(*fiber.Ctx) error
	GetClassroom(*fiber.Ctx) error
	UpdateClassroom(*fiber.Ctx) error
	ArchiveClassroom(*fiber.Ctx) error

	GetClassroomTemplates(*fiber.Ctx) error

	GetClassroomAssignments(*fiber.Ctx) error
	CreateAssignment(*fiber.Ctx) error
	ClassroomAssignmentMiddleware(*fiber.Ctx) error
	GetClassroomAssignment(*fiber.Ctx) error
	UpdateAssignment(*fiber.Ctx) error

	GetGradingRubrics(c *fiber.Ctx) (err error)
	UpdateGradingRubrics(c *fiber.Ctx) (err error)

	GetClassroomAssignmentProjects(*fiber.Ctx) error
	InviteToAssignment(*fiber.Ctx) error
	ClassroomAssignmentProjectMiddleware(*fiber.Ctx) error
	GetClassroomAssignmentProject(*fiber.Ctx) error

	GetGradingResults(c *fiber.Ctx) (err error)
	UpdateGradingResults(c *fiber.Ctx) (err error)

	StartAutoGrading(c *fiber.Ctx) (err error)
	StartAutoGradingForProject(c *fiber.Ctx) (err error)

	GetClassroomReport(c *fiber.Ctx) (err error)
	GetClassroomAssignmentReport(c *fiber.Ctx) (err error)
	GetClassroomTeamReport(c *fiber.Ctx) (err error)

	GetClassroomProjects(*fiber.Ctx) error
	AcceptAssignment(*fiber.Ctx) error
	ClassroomProjectMiddleware(*fiber.Ctx) error
	GetClassroomProject(*fiber.Ctx) error

	GetClassroomInvitation(*fiber.Ctx) error
	JoinClassroom(*fiber.Ctx) error
	GetClassroomInvitations(*fiber.Ctx) error
	InviteToClassroom(*fiber.Ctx) error
	RevokeClassroomInvitation(*fiber.Ctx) error

	GetClassroomMembers(*fiber.Ctx) error
	ClassroomMemberMiddleware(*fiber.Ctx) error
	GetClassroomMember(*fiber.Ctx) error
	UpdateMemberTeam(*fiber.Ctx) error
	UpdateMemberRole(*fiber.Ctx) error

	GetClassroomRunners(c *fiber.Ctx) error
	GetClassroomRunnersAreAvailable(c *fiber.Ctx) error

	GetClassroomTeams(*fiber.Ctx) error
	CreateTeam(*fiber.Ctx) error
	JoinTeam(*fiber.Ctx) error
	ClassroomTeamMiddleware(*fiber.Ctx) error
	GetClassroomTeam(*fiber.Ctx) error
	UpdateTeam(*fiber.Ctx) error

	GetClassroomTeamMembers(*fiber.Ctx) error
	ClassroomTeamMemberMiddleware(*fiber.Ctx) error
	GetClassroomTeamMember(*fiber.Ctx) error
	RemoveMemberFromTeam(*fiber.Ctx) error

	GetClassroomTeamProjects(*fiber.Ctx) error
	ClassroomTeamProjectMiddleware(*fiber.Ctx) error
	GetClassroomTeamProject(*fiber.Ctx) error
	GetGitlabInfo(*fiber.Ctx) error
}
