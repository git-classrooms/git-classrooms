package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type assignmentGradingResponse struct {
	GradingJUnitAutoGradingActive bool `json:"gradingJUnitAutoGradingActive"`

	GradingManualRubrics []*database.ManualGradingRubric `json:"gradingManualRubrics"`
} //@Name AssignmentGradingResponse

// @Summary		GetGradingRubrics
// @Description	GetGradingRubrics
// @Id				GetGradingRubrics
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{object}	api.AssignmentGradingResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading [get]
func (ctrl *DefaultController) GetGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()

	queryManualGradingRubric := query.ManualGradingRubric
	manualGradingRubrics, err := queryManualGradingRubric.
		WithContext(c.Context()).
		Where(queryManualGradingRubric.AssignmentID.Eq(assignment.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(assignmentGradingResponse{
		GradingJUnitAutoGradingActive: assignment.GradingJUnitAutoGradingActive,
		GradingManualRubrics:          manualGradingRubrics,
	})
}
