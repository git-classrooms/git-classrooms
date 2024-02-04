package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetMe(*fiber.Ctx) error

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
