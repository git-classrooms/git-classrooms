package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) JoinAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()
	team := classroom.Team
	assignmentProject := ctx.GetJoinedClassroomAssignment()

	if assignmentProject.AssignmentAccepted {
		return c.SendStatus(fiber.StatusNoContent) // You or your teammate have already accepted the assignment
	}

	repo := context.Get(c).GetGitlabRepository()

	if err := repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if template repository still exists
	_, err := repo.GetProjectById(assignmentProject.Assignment.TemplateProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	project, err := repo.ForkProject(assignmentProject.Assignment.TemplateProjectID, gitlabModel.Private, team.GroupID, assignmentProject.Assignment.Name, assignmentProject.Assignment.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	assignmentProject.ProjectID = project.ID
	assignmentProject.AssignmentAccepted = true
	queryAssignmentProjects := query.AssignmentProjects
	err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject)
	if err != nil {
		if err := repo.DeleteProject(project.ID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
