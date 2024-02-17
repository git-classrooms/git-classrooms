package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) GetClassroomAssignment(c *fiber.Ctx) error {
	param := Params{}
	err := c.ParamsParser(&param)
	if err != nil || param.ClassroomID == nil || param.AssignmentID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryAssignment := query.Assignment
	assignment, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ClassroomID.Eq(param.ClassroomID)).
		Where(queryAssignment.ID.Eq(param.AssignmentID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(assignment)
}
