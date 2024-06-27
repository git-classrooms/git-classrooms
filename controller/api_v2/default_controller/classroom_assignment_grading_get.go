package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type gradingManualRubricResponse struct {
	ID          *uuid.UUID `json:"id" validate:"optional"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	MaxScore    int        `json:"maxScore"`
} //@Name GradingManualRubricResponse

// @Summary		GetGradingRubrics
// @Description	GetGradingRubrics
// @Id				GetGradingRubrics
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path	string								true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string								true	"Assignment ID"	Format(uuid)
// @Success		200				{array}	api.gradingManualRubricResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading [get]
func (ctrl *DefaultController) GetGradingRubrics(c *fiber.Ctx) (err error) {
	return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
}
