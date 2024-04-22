package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func ownedClassroomAssignmentProjectQuery(assignmentId uuid.UUID, c *fiber.Ctx) query.IAssignmentProjectsDo {
	queryAssignmentProject := query.AssignmentProjects
	return queryAssignmentProject.
		WithContext(c.Context()).
		Where(queryAssignmentProject.AssignmentID.Eq(assignmentId))
}

func (ctrl *DefaultController) OwnedClassroomAssignmentProjectMiddleware(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.AssignmentID == nil || param.AssignmentProjectID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryAssignmentProject := query.AssignmentProjects
	assignmentProject, err := ownedClassroomAssignmentProjectQuery(*param.AssignmentID, c).
		Where(queryAssignmentProject.ID.Eq(*param.AssignmentProjectID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetOwnedClassroomAssignmentProject(assignmentProject)
	ctx.SetGitlabProjectID(assignmentProject.ProjectID)
	return ctx.Next()
}
