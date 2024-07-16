package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type startAutoGradingRequest struct {
	JUnitAutoGrading *bool `json:"jUnitAutoGrading"`
} //@Name StartAutoGradingRequest

func (r startAutoGradingRequest) isValid() bool {
	return r.JUnitAutoGrading != nil
}

// @Summary		StartAutoGrading
// @Description	StartAutoGrading
// @Id				StartAutoGrading
// @Tags			grading
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string						true	"Assignment ID"	Format(uuid)
// @Param			assignmentInfo	body	api.startAutoGradingRequest	true	"Grading Update Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading/auto [post]
func (ctrl *DefaultController) StartAutoGrading(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	repo := ctx.GetGitlabRepository()
	assignment := ctx.GetAssignment()

	var requestBody startAutoGradingRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	queryAssignmentProjects := query.AssignmentProjects
	projects, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Where(queryAssignmentProjects.AssignmentID.Eq(assignment.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if *requestBody.JUnitAutoGrading {
		if !assignment.GradingJUnitAutoGradingActive {
			return fiber.NewError(fiber.StatusBadRequest, "JUnit Auto Grading is not active")
		}

		for _, project := range projects {
			if err = InsertJUnitTestResultForProject(c.Context(), repo, project); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func InsertJUnitTestResultForProject(ctx context.Context, repo gitlab.Repository, project *database.AssignmentProjects) error {
	report, err := repo.GetProjectLatestPipelineTestReportSummary(project.ProjectID, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return query.Q.Transaction(func(tx *query.Query) (err error) {
		if project.GradingJUnitTestResultID != nil {
			_, err = tx.JUnitTestResult.
				WithContext(ctx).
				Where(tx.JUnitTestResult.ID.Eq(*project.GradingJUnitTestResultID)).
				Delete()
			if err != nil {
				return err
			}
		}

		result := database.JUnitTestResult{
			TotalTime:    report.TotalTime,
			TotalCount:   report.TotalCount,
			SuccessCount: report.SuccessCount,
			FailedCount:  report.FailedCount,
			SkippedCount: report.SkippedCount,
			ErrorCount:   report.ErrorCount,
		}

		err = tx.JUnitTestResult.WithContext(ctx).Create(&result)
		if err != nil {
			return err
		}

		project.GradingJUnitTestResultID = &result.ID
		return tx.AssignmentProjects.WithContext(ctx).Save(project)
	})
}
