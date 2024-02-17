package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getMeClassroomAssignmentResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetMeClassroomAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetClassroom()

	assignmentId, err := uuid.Parse(c.Params("assignmentId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no assignment specified")
	}

	queryAssignment := query.Assignment
	queryAssignmentProjects := query.AssignmentProjects
	assignmentProject, err := queryAssignmentProjects.WithContext(c.Context()).
		Join(queryAssignment, queryAssignment.ID.EqCol(queryAssignmentProjects.AssignmentID)).
		Where(queryAssignment.ClassroomID.Eq(classroom.ClassroomID)).
		Where(queryAssignment.ID.Eq(assignmentId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	repo := ctx.GetGitlabRepository()
	webURL := ""
	if assignmentProject.AssignmentAccepted {
		projectFromGitLab, err := repo.GetProjectById(assignmentProject.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		webURL = projectFromGitLab.WebUrl
	}
	response := &getMeClassroomAssignmentResponse{
		AssignmentProjects: *assignmentProject,
		ProjectPath:        webURL,
	}

	return c.JSON(response)
}
