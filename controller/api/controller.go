package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetMe(*fiber.Ctx) error
	GetMeGitlab(*fiber.Ctx) error

	RedirectGroupGitlab(*fiber.Ctx) error
	RedirectUserGitlab(*fiber.Ctx) error
	RedirectProjectGitlab(*fiber.Ctx) error

	GetOwnedClassrooms(*fiber.Ctx) error
	OwnedClassroomMiddleware(*fiber.Ctx) error
	GetOwnedClassroom(*fiber.Ctx) error
	PutOwnedClassroom(*fiber.Ctx) error

	GetOwnedClassroomAssignments(*fiber.Ctx) error
	OwnedClassroomAssignmentMiddleware(*fiber.Ctx) error
	GetOwnedClassroomAssignment(*fiber.Ctx) error

	GetOwnedClassroomAssignmentProjects(*fiber.Ctx) error
	OwnedClassroomAssignmentProjectMiddleware(*fiber.Ctx) error
	GetOwnedClassroomAssignmentProject(*fiber.Ctx) error

	InviteToAssignmentProject(*fiber.Ctx) error

	GetOwnedClassroomMembers(*fiber.Ctx) error
	OwnedClassroomMemberMiddleware(*fiber.Ctx) error
	GetOwnedClassroomMember(*fiber.Ctx) error
	ChangeOwnedClassroomMember(*fiber.Ctx) error

	GetOwnedClassroomInvitations(*fiber.Ctx) error

	GetOwnedClassroomTemplates(*fiber.Ctx) error

	GetOwnedClassroomTeams(*fiber.Ctx) error
	OwnedClassroomTeamMiddleware(*fiber.Ctx) error
	GetOwnedClassroomTeam(*fiber.Ctx) error
	CreateOwnedClassroomTeam(*fiber.Ctx) error

	GetOwnedClassroomTeamMembers(*fiber.Ctx) error
	OwnedClassroomTeamMemberMiddleware(*fiber.Ctx) error
	RemoveMemberFromTeam(*fiber.Ctx) error

	GetOwnedClassroomTeamProjects(*fiber.Ctx) error

	GetJoinedClassrooms(*fiber.Ctx) error
	JoinedClassroomMiddleware(*fiber.Ctx) error
	GetJoinedClassroom(*fiber.Ctx) error

	GetJoinedClassroomAssignments(*fiber.Ctx) error
	JoinedClassroomAssignmentMiddleware(*fiber.Ctx) error
	GetJoinedClassroomAssignment(*fiber.Ctx) error

	GetJoinedClassroomTeams(*fiber.Ctx) error
	CreateJoinedClassroomTeam(*fiber.Ctx) error
	JoinedClassroomTeamMiddleware(*fiber.Ctx) error
	GetJoinedClassroomTeam(*fiber.Ctx) error
	JoinJoinedClassroomTeam(*fiber.Ctx) error

	JoinClassroom(*fiber.Ctx) error
	GetInvitationInfo(*fiber.Ctx) error
	AcceptAssignment(*fiber.Ctx) error

	CreateClassroom(*fiber.Ctx) error
	CreateAssignment(*fiber.Ctx) error

	InviteToClassroom(*fiber.Ctx) error
}
