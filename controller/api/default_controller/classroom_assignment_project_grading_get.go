package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type projectGradingResponse struct {
	GradingJUnitTestResult *database.JUnitTestResult       `json:"gradingJUnitTestResult"`
	GradingManualResults   []*database.ManualGradingResult `json:"gradingManualResults"`
} //@Name AssignmentGradingResponse

// @Summary		GetGradingResults
// @Description	GetGradingResults
// @Tags			grading
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Param			projectId		path		string	true	"Project ID"	Format(uuid)
// @Success		200				{array}		ManualGradingResult
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/projects/{projectId}/grading [get]
// @Router			/api/v2/classrooms/{classroomId}/projects/{projectId}/grading [get]
func (ctrl *DefaultController) GetGradingResults(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	project := ctx.GetAssignmentProject()

	queryManualGradingResult := query.ManualGradingResult
	results, err := queryManualGradingResult.
		WithContext(c.Context()).
		Preload(queryManualGradingResult.Rubric).
		Where(queryManualGradingResult.AssignmentProjectID.Eq(project.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(projectGradingResponse{
		GradingJUnitTestResult: project.GradingJUnitTestResult,
		GradingManualResults:   results,
	})
}
