package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func teamProjectQuery(c *fiber.Ctx, teamID uuid.UUID) query.IAssignmentProjectsDo {
	queryAssignmentProject := query.AssignmentProjects
	return queryAssignmentProject.
		WithContext(c.Context()).
		Where(queryAssignmentProject.TeamID.Eq(teamID))
}

func (ctrl *DefaultController) ClassroomTeamProjectMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.TeamID == nil || params.AssignmentProjectID == nil {
		return fiber.ErrBadRequest
	}

	assignmentProject, err := teamProjectQuery(c, *params.TeamID).
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
