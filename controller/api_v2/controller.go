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
	// CreateClassroom(*fiber.Ctx) error
	ClassroomMiddleware(*fiber.Ctx) error // Implemented
	GetClassroom(*fiber.Ctx) error        // Implemented

	GetClassroomTemplates(*fiber.Ctx) error

	GetClassroomAssignments(*fiber.Ctx) error // Implemented
	// CreateAssignment(*fiber.Ctx) error
	ClassroomAssignmentMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomAssignment(*fiber.Ctx) error        // Implemented

	GetClassroomAssignmentProjects(*fiber.Ctx) error // Implemented
	// InviteToAssignment(*fiber.Ctx) error
	ClassroomAssignmentProjectMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomAssignmentProject(*fiber.Ctx) error        // Implemented

	GetClassroomProjects(*fiber.Ctx) error // Implemented
	// AcceptAssignmentProject(*fiber.Ctx) error
	ClassroomProjectMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomProject(*fiber.Ctx) error        // Implemented

	GetClassroomInvitation(*fiber.Ctx) error // Implemented
	// JoinClassroom(*fiber.Ctx) error
	GetClassroomInvitations(*fiber.Ctx) error // Implemented
	// InviteToClassroom(*fiber.Ctx) error

	GetClassroomMembers(*fiber.Ctx) error       // Implemented
	ClassroomMemberMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomMember(*fiber.Ctx) error        // Implemented

	GetClassroomTeams(*fiber.Ctx) error // Implemented
	// CreateTeam(*fiber.Ctx) error
	// JoinTeam(*fiber.Ctx) error
	ClassroomTeamMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomTeam(*fiber.Ctx) error        // Implemented

	GetClassroomTeamMembers(*fiber.Ctx) error       // Implemented
	ClassroomTeamMemberMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomTeamMember(*fiber.Ctx) error        // Implemented

	GetClassroomTeamProjects(*fiber.Ctx) error       // Implemented
	ClassroomTeamProjectMiddleware(*fiber.Ctx) error // Implemented
	GetClassroomTeamProject(*fiber.Ctx) error        // Implemented
}
