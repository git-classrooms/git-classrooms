package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetMe(*fiber.Ctx) error

	GetOwnedClassrooms(*fiber.Ctx) error
	OwnedClassroomMiddleware(*fiber.Ctx) error
	GetOwnedClassroom(*fiber.Ctx) error

	GetOwnedClassroomAssignments(*fiber.Ctx) error
	OwnedClassroomAssignmentMiddleware(*fiber.Ctx) error
	GetOwnedClassroomAssignment(*fiber.Ctx) error

	GetOwnedClassroomAssignmentProjects(*fiber.Ctx) error
	InviteToAssignmentProject(*fiber.Ctx) error

	GetOwnedClassroomMembers(*fiber.Ctx) error

	GetOwnedClassroomInvitations(*fiber.Ctx) error

	GetOwnedClassroomTemplates(*fiber.Ctx) error

	GetJoinedClassrooms(*fiber.Ctx) error
	JoinedClassroomMiddleware(*fiber.Ctx) error
	GetJoinedClassroom(*fiber.Ctx) error

	GetJoinedClassroomAssignments(*fiber.Ctx) error
	JoinedClassroomAssignmentMiddleware(*fiber.Ctx) error
	GetJoinedClassroomAssignment(*fiber.Ctx) error

	JoinClassroomNew(*fiber.Ctx) error
	JoinAssignmentNew(*fiber.Ctx) error

	GetMeClassrooms(*fiber.Ctx) error
	GetMeClassroomMiddleware(*fiber.Ctx) error
	GetMeClassroom(*fiber.Ctx) error
	GetMeClassroomTemplates(*fiber.Ctx) error
	GetMeClassroomInvitations(*fiber.Ctx) error
	GetMeClassroomMember(*fiber.Ctx) error
	GetMeClassroomMemberAssignment(*fiber.Ctx) error
	GetMeClassroomMemberAssignments(*fiber.Ctx) error
	GetMeClassroomMembers(*fiber.Ctx) error
	GetMeClassroomAssignments(*fiber.Ctx) error
	GetMeClassroomAssignment(*fiber.Ctx) error

	CreateClassroom(*fiber.Ctx) error
	GetClassroomAssignments(ctx *fiber.Ctx) error
	GetClassroomAssignment(ctx *fiber.Ctx) error
	GetClassroomAssignmentProjects(*fiber.Ctx) error
	CreateAssignment(*fiber.Ctx) error
	JoinClassroom(*fiber.Ctx) error
	JoinAssignment(*fiber.Ctx) error
	InviteToClassroom(*fiber.Ctx) error
	InviteToAssignment(*fiber.Ctx) error
}
