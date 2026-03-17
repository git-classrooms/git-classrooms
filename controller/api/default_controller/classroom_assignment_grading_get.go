package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type GetAssignmentGradingRubricsResponse struct {
	AssignmentDateID uuid.UUID                       `json:"assignmentDateId"`
	Rubrics          []*database.ManualGradingRubric `json:"rubrics"`
} //@Name GetAssignmentGradingRubricsResponse

// @Summary		GetAssignmentGradingRubrics
// @Description	GetAssignmentGradingRubrics
// @Id				GetAssignmentGradingRubrics
// @Tags			grading
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{array}		api.GetAssignmentGradingRubricsResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/grading [get]
func (ctrl *DefaultController) GetAssignmentGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()

	response := utils.Map(assignment.AssignmentDates, func(ad *database.AssignmentDate) GetAssignmentGradingRubricsResponse {
		return GetAssignmentGradingRubricsResponse{
			AssignmentDateID: ad.ID,
			Rubrics:          ad.GradingManualRubrics,
		}
	})

	return c.JSON(response)
}
