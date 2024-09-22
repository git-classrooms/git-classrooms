package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gen/field"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func classroomProjectQuery(c *fiber.Ctx, classroomID uuid.UUID, teamID uuid.UUID) query.IAssignmentProjectsDo {
	queryAssignment := query.Assignment
	queryAssignmentProjects := query.AssignmentProjects
	return queryAssignmentProjects.
		WithContext(c.Context()).
		Preload(queryAssignmentProjects.Assignment).
		Preload(queryAssignmentProjects.Team).
		Preload(queryAssignmentProjects.GradingManualResults).
		Preload(queryAssignmentProjects.GradingManualResults.Rubric).
		Preload(field.NewRelation("Team.Member", "")).
		Join(queryAssignment, queryAssignment.ID.EqCol(queryAssignmentProjects.AssignmentID)).
		Where(queryAssignment.ClassroomID.Eq(classroomID)).
		Where(queryAssignmentProjects.TeamID.Eq(teamID))
}

func (ctrl *DefaultController) ClassroomProjectMiddleware(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	if classroom.TeamID == nil {
		return fiber.ErrForbidden
	}

	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.AssignmentProjectID == nil {
		return fiber.ErrBadRequest
	}

	project, err := classroomProjectQuery(c, *params.ClassroomID, *classroom.TeamID).
		Where(query.AssignmentProjects.ID.Eq(*params.AssignmentProjectID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetAssignmentProject(project)
	ctx.SetGitlabProjectID(project.ProjectID)

	return c.Next()
}
