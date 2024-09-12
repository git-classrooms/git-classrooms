package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetAssignmentGradingRubrics
// @Description	GetAssignmentGradingRubrics
// @Id				GetAssignmentGradingRubrics
// @Tags			grading
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{array}		ManualGradingRubric
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading [get]
func (ctrl *DefaultController) GetAssignmentGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()

	return c.JSON(assignment.GradingManualRubrics)
}
