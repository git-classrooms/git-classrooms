package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetGradingRubrics
// @Description	GetGradingRubrics
// @Id				GetGradingRubrics
// @Tags			grading
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		ManualGradingRubric
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/grading [get]
func (ctrl *DefaultController) GetGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	queryManualGradingRubric := query.ManualGradingRubric
	manualGradingRubrics, err := queryManualGradingRubric.
		WithContext(c.Context()).
		Where(queryManualGradingRubric.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(manualGradingRubrics)
}
