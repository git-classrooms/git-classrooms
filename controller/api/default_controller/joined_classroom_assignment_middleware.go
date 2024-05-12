package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func joinedClassroomAssignmentQuery(classroomID uuid.UUID, teamID uuid.UUID, c *fiber.Ctx) query.IAssignmentProjectsDo {
	queryAssignment := query.Assignment
	queryAssignmentProjects := query.AssignmentProjects
	return queryAssignmentProjects.
		WithContext(c.Context()).
		Preload(queryAssignmentProjects.Assignment).
		Preload(queryAssignmentProjects.Team).
		Preload(queryAssignmentProjects.Team.Member).
		Join(queryAssignment, queryAssignment.ID.EqCol(queryAssignmentProjects.AssignmentID)).
		Where(queryAssignment.ClassroomID.Eq(classroomID)).
		Where(queryAssignmentProjects.TeamID.Eq(teamID))
}

func (ctrl *DefaultController) JoinedClassroomAssignmentMiddleware(c *fiber.Ctx) error {
	classroom := context.Get(c).GetJoinedClassroom()

	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.AssignmentID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignment, err := joinedClassroomAssignmentQuery(*param.ClassroomID, *classroom.TeamID, c).
		Where(query.AssignmentProjects.AssignmentID.Eq(*param.AssignmentID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetJoinedClassroomAssignment(assignment)
	return ctx.Next()
}
