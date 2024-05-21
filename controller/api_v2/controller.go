package apiV2

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

type ValidateUserFunc func(user database.UserClassrooms) bool

type Controller interface {
	ValidateUserMiddleware(ValidateUserFunc) fiber.Handler    // Implemented
	RoleMiddleware(validRoles ...database.Role) fiber.Handler // Implemented
	ViewableClassroomMiddleware() fiber.Handler               // Implemented

	RedirectUserGitlab(*fiber.Ctx) error    // Implemented
	RedirectGroupGitlab(*fiber.Ctx) error   // Implemented
	RedirectProjectGitlab(*fiber.Ctx) error // Implemented

	GetMe(*fiber.Ctx) error       // Implemented
	GetMeGitlab(*fiber.Ctx) error // Implemented

	GetClassrooms(*fiber.Ctx) error // Implemented
	CreateClassroom(*fiber.Ctx) error
	ClassroomMiddleware(*fiber.Ctx) error // Implemented
	GetClassroom(*fiber.Ctx) error        // Implemented

	GetClassroomTemplates(*fiber.Ctx) error

	GetClassroomAssignments(*fiber.Ctx) error // Implemented
	CreateAssignment(*fiber.Ctx) error
	ClassroomAssignmentMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomAssignment(*fiber.Ctx) error        // Implemented

	GetClassroomAssignmentProjects(*fiber.Ctx) error
	InviteToAssignment(*fiber.Ctx) error
	ClassroomAssignmentProjectMiddleware(*fiber.Ctx) error
	GetClassroomAssignmentProject(*fiber.Ctx) error

	GetClassroomProjects(*fiber.Ctx) error // Implemented
	AcceptAssignmentProject(*fiber.Ctx) error
	ClassroomProjectMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomProject(*fiber.Ctx) error        // Implemented

	GetClassroomInvitation(*fiber.Ctx) error
	JoinClassroom(*fiber.Ctx) error
	GetClassroomInvitations(*fiber.Ctx) error
	InviteToClassroom(*fiber.Ctx) error
	ClassroomInvitationMiddleware(*fiber.Ctx) error

	GetClassroomMembers(*fiber.Ctx) error
	ClassroomMemberMiddleware(*fiber.Ctx) error
	GetClassroomMember(*fiber.Ctx) error

	GetClassroomTeams(*fiber.Ctx) error
	CreateTeam(*fiber.Ctx) error
	JoinTeam(*fiber.Ctx) error
	ClassroomTeamMiddleware(*fiber.Ctx) error
	GetClassroomTeam(*fiber.Ctx) error

	GetClassroomTeamMembers(*fiber.Ctx) error
	ClassroomTeamMemberMiddleware(*fiber.Ctx) error
	GetClassroomTeamMember(*fiber.Ctx) error

	GetClassroomTeamProjects(*fiber.Ctx) error
	ClassroomTeamProjectMiddleware(*fiber.Ctx) error
	GetClassroomTeamProject(*fiber.Ctx) error
}
