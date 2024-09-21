package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

// @Summary		StartAutoGradingForProject
// @Description	StartAutoGradingForProject
// @Id				StartAutoGradingForProject
// @Tags			grading
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string						true	"Assignment ID"	Format(uuid)
// @Param			projectId		path	string						true	"Project ID"	Format(uuid)
// @Param			assignmentInfo	body	api.startAutoGradingRequest	true	"Grading Update Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/projects/{projectId}/grading/auto [post]
func (ctrl *DefaultController) StartAutoGradingForProject(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	classroom := ctx.GetUserClassroom()
	repo := ctx.GetGitlabRepository()
	if err := repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	project := ctx.GetAssignmentProject()

	var requestBody startAutoGradingRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	if *requestBody.JUnitAutoGrading {
		if !project.Assignment.GradingJUnitAutoGradingActive {
			return fiber.NewError(fiber.StatusBadRequest, "JUnit Auto Grading is not active")
		}

		report, err := repo.GetProjectLatestPipelineTestReportSummary(project.ProjectID, nil)
		if err != nil {
			var gitlabError *model.GitLabError
			if errors.As(err, &gitlabError) {
				if gitlabError.Response.StatusCode == http.StatusForbidden {
					return fiber.NewError(fiber.StatusNotFound, "No executed pipeline yet available on the main branch")
				}
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		project.GradingJUnitTestResult = &database.JUnitTestResult{TestReport: *report}

		if err := query.AssignmentProjects.WithContext(c.Context()).Save(project); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusAccepted)
}
