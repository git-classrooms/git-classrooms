package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"log"
)

func (ctrl *DefaultController) JoinAssignmentNew(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()
	userID := ctx.GetUserID()
	// TODO: assignnment := ctx.GetJoinedClassroomAssignment()

	assignmentId, err := uuid.Parse(c.Params("assignmentId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: Delete this - Check if assignment is valid
	queryAssignment := query.Assignment
	assignment, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ID.Eq(assignmentId)).
		Where(queryAssignment.ClassroomID.Eq(classroom.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if invitation is valid
	queryAssignmentProjects := query.AssignmentProjects
	assignmentProject, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Preload(queryAssignmentProjects.Assignment).
		Where(queryAssignmentProjects.AssignmentID.Eq(assignmentId)).
		Where(queryAssignmentProjects.UserID.Eq(userID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if assignmentProject.AssignmentAccepted {
		return fiber.NewError(fiber.StatusBadRequest, "You have already joined this assignment")
	}

	repo := context.Get(c).GetGitlabRepository()
	// Check if template repository still exists
	_, err = repo.GetProjectById(assignment.TemplateProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := repo.GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	assignmentName := fmt.Sprintf("%s-%s", assignment.Name, user.Username)
	log.Println(assignmentName)

	project, err := repo.ForkProject(assignment.TemplateProjectID, gitlabModel.Private, classroom.Classroom.GroupID, assignmentName, assignment.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	assignmentProject.ProjectID = project.ID
	assignmentProject.AssignmentAccepted = true
	err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject)
	if err != nil {
		err := repo.DeleteProject(project.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
