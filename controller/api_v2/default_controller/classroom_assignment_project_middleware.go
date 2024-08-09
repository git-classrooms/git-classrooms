package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func assignmentProjectQuery(c *fiber.Ctx, assignmentID uuid.UUID) query.IAssignmentProjectsDo {
	queryAssignmentProject := query.AssignmentProjects
	return queryAssignmentProject.
		WithContext(c.Context()).
		Preload(queryAssignmentProject.Team).
		Preload(queryAssignmentProject.GradingManualResults).
		Preload(queryAssignmentProject.GradingManualResults.Rubric).
		Preload(queryAssignmentProject.GradingJUnitTestResult).
		Where(queryAssignmentProject.AssignmentID.Eq(assignmentID))
}

func (ctrl *DefaultController) ClassroomAssignmentProjectMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.AssignmentID == nil || params.AssignmentProjectID == nil {
		return fiber.ErrBadRequest
	}

	assignmentProject, err := assignmentProjectQuery(c, *params.AssignmentID).
		Where(query.AssignmentProjects.ID.Eq(*params.AssignmentProjectID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetAssignmentProject(assignmentProject)
	ctx.SetGitlabProjectID(assignmentProject.ProjectID)

	return c.Next()
}
