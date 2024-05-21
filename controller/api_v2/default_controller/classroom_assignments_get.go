package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomAssignments
// @Description	GetClassroomAssignments
// @Id				GetClassroomAssignments
// @Tags			assignment
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		Assignment
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments [get]
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
