package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetOwnedClassroomAssignment
// @Description	GetOwnedClassroomAssignment
// @Id				GetOwnedClassroomAssignment
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{object}	database.Assignment
// @Failure		400				{object}	httputil.HTTPError
// @Failure		401				{object}	httputil.HTTPError
// @Failure		404				{object}	httputil.HTTPError
// @Failure		500				{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/assignments/{assignmentId} [get]
func (ctrl *DefaultController) GetOwnedClassroomAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetOwnedClassroomAssignment()
	return c.JSON(assignment)
}
