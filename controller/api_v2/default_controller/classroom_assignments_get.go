package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomAssignments(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	assignments, err := classroomAssignmentQuery(c, classroom.ClassroomID).
		Find()
	if err != nil {
		return err
	}

	return c.JSON(assignments)
}
