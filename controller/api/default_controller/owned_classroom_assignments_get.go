package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetOwnedClassroomAssignments
// @Description	GetOwnedClassroomAssignments
// @Id				GetOwnedClassroomAssignments
// @Tags			assignment
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		Assignment
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/classrooms/owned/{classroomId}/assignments [get]
func (ctrl *DefaultController) GetOwnedClassroomAssignments(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	queryAssignment := query.Assignment
	assignments, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ClassroomID.Eq(classroom.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(assignments)
}
