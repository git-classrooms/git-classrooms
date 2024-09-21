package api

import (
	"database/sql/driver"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateAssignmentRubricsRequest struct {
	RubricIDs []uuid.UUID `json:"rubricIds"`
} //@Name UpdateAssignmentRubricsRequest

func (r updateAssignmentRubricsRequest) isValid() bool {
	return r.RubricIDs != nil
}

// @Summary		UpdateAssignmentGradingRubrics
// @Description	UpdateAssignmentGradingRubrics
// @Id				UpdateAssignmentGradingRubrics
// @Tags			grading
// @Accept			json
// @Param			classroomId				path	string								true	"Classroom ID"	Format(uuid)
// @Param			assignmentId			path	string								true	"Assignment ID"	Format(uuid)
// @Param			assignmentGradingInfo	body	api.updateAssignmentRubricsRequest	true	"Assignment Grading Update Info"
// @Param			X-Csrf-Token			header	string								true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading [put]
func (ctrl *DefaultController) UpdateAssignmentGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()

	var requestBody updateAssignmentRubricsRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	ids := utils.Map(requestBody.RubricIDs, func(e uuid.UUID) driver.Valuer { return e })

	queryManualGradingRubric := query.ManualGradingRubric
	rubrics, err := queryManualGradingRubric.
		WithContext(c.Context()).
		Where(queryManualGradingRubric.ClassroomID.Eq(assignment.ClassroomID)).
		Where(queryManualGradingRubric.ID.In(ids...)).Find()
	if err != nil {
		return err
	}

	if len(rubrics) != len(ids) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid IDs")
	}

	if err := query.Assignment.GradingManualRubrics.Model(assignment).Replace(rubrics...); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryManualGradingResult := query.ManualGradingResult
	toDeleteResults, err := queryManualGradingResult.
		WithContext(c.Context()).
		Join(queryManualGradingRubric, queryManualGradingResult.RubricID.EqCol(queryManualGradingRubric.ID)).
		Where(queryManualGradingRubric.ClassroomID.Eq(assignment.ClassroomID)).
		Not(queryManualGradingResult.RubricID.In(ids...)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ids = utils.Map(toDeleteResults, func(e *database.ManualGradingResult) driver.Valuer { return e.ID })

	if _, err := queryManualGradingResult.WithContext(c.Context()).Where(queryManualGradingResult.ID.In(ids...)).Delete(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
