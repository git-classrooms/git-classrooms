package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

type getGetClassroomAssignmentProjectsResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetClassroomAssignmentProjects(c *fiber.Ctx) error {
	classroomId, err := uuid.Parse(c.Params("classroomId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignmentId, err := uuid.Parse(c.Params("assignmentId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryAssignment := query.Assignment
	assignment, err := queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Projects).
		Preload(queryAssignment.Projects.User).
		Where(queryAssignment.ClassroomID.Eq(classroomId)).
		Where(queryAssignment.ID.Eq(assignmentId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	responses := make([]getGetClassroomAssignmentProjectsResponse, len(assignment.Projects))
	for i, project := range assignment.Projects {
		repo := context.Get(c).GetGitlabRepository()
		webURL := ""
		if project.AssignmentAccepted {
			projectFromGitLab, err := repo.GetProjectById(project.ProjectID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			webURL = projectFromGitLab.WebUrl
		}
		responses[i] = getGetClassroomAssignmentProjectsResponse{
			AssignmentProjects: *project,
			ProjectPath:        webURL,
		}
	}

	return c.JSON(responses)
}
