package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func ownedClassroomAssignmentQuery(classroomId uuid.UUID, c *fiber.Ctx) query.IAssignmentDo {
	queryAssignment := query.Assignment
	return queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Projects).
		Preload(queryAssignment.Projects.User).
		Where(queryAssignment.ClassroomID.Eq(classroomId))
}

func (ctrl *DefaultController) OwnedClassroomAssignmentMiddleware(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.AssignmentID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignment, err := ownedClassroomAssignmentQuery(*param.ClassroomID, c).
		Where(query.Assignment.ID.Eq(*param.AssignmentID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetOwnedClassroomAssignment(assignment)
	return ctx.Next()
}
