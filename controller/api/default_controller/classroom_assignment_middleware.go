package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func classroomAssignmentQuery(c *fiber.Ctx, classroomID uuid.UUID) query.IAssignmentDo {
	queryAssignment := query.Assignment

	return queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.GradingManualRubrics).
		Preload(queryAssignment.JUnitTests).
		Where(queryAssignment.ClassroomID.Eq(classroomID))
}

func (*DefaultController) ClassroomAssignmentMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.AssignmentID == nil {
		return fiber.ErrBadRequest
	}

	assignment, err := classroomAssignmentQuery(c, *params.ClassroomID).
		Where(query.Assignment.ID.Eq(*params.AssignmentID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetAssignment(assignment)
	return ctx.Next()
}
