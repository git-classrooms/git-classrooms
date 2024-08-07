package api

import (
	"github.com/gofiber/fiber/v2"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
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
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/projects/{projectId}/grading/auto [post]
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

		if err = InsertJUnitTestResultForProject(c.Context(), repo, project); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusAccepted)
}
