package api

import (
	"database/sql/driver"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type gradingManualResultRequest struct {
	RubricID *uuid.UUID `json:"id"`
	Score    *int       `json:"score"`
	Feedback *string    `json:"feedback" validate:"optional"`
} //@Name GradingManualResultRequest

func (r gradingManualResultRequest) isValid() bool {
	return r.RubricID != nil && r.Score != nil
}

func resultRequestIsValid(r gradingManualResultRequest) bool {
	return r.isValid()
}

type updateProjectGradingRequest struct {
	GradingManualResults []gradingManualResultRequest `json:"gradingManualRubrics"`
} //@Name UpdateAssignmentGradingRequest

func (r updateProjectGradingRequest) isValid() bool {
	return utils.All(r.GradingManualResults, resultRequestIsValid)
}

// @Summary		UpdateGradingResults
// @Description	UpdateGradingResults
// @Id				UpdateGradingResults
// @Tags			grading
// @Accept			json
// @Param			classroomId		path	string							true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string							true	"Assignment ID"	Format(uuid)
// @Param			projectId		path	string							true	"Project ID"	Format(uuid)
// @Param			assignmentInfo	body	api.updateProjectGradingRequest	true	"Grading Update Info"
// @Param			X-Csrf-Token	header	string							true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/projects/{projectId}/grading [put]
func (ctrl *DefaultController) UpdateGradingResults(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	project := ctx.GetAssignmentProject()

	var requestBody updateProjectGradingRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	rubricIDs := utils.Map(requestBody.GradingManualResults, func(e gradingManualResultRequest) driver.Valuer { return *e.RubricID })

	if len(assignment.GradingManualRubrics) != len(rubricIDs) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes not enough rubrics")
	}

	if !utils.All(rubricIDs, func(e driver.Valuer) bool {
		return slices.ContainsFunc(assignment.GradingManualRubrics, func(rubric *database.ManualGradingRubric) bool { return rubric.ID == e.(uuid.UUID) })
	}) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid IDs")
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		queryManualGradingResult := tx.ManualGradingResult
		if _, err := queryManualGradingResult.
			WithContext(c.Context()).
			Where(queryManualGradingResult.AssignmentProjectID.Eq(project.ID)).
			Delete(); err != nil {
			return err
		}

		results := utils.Map(requestBody.GradingManualResults, func(e gradingManualResultRequest) *database.ManualGradingResult {
			return &database.ManualGradingResult{
				AssignmentProjectID: project.ID,
				RubricID:            *e.RubricID,
				Score:               *e.Score,
				Feedback:            e.Feedback,
			}
		})

		if err = queryManualGradingResult.
			WithContext(c.Context()).
			Save(results...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
