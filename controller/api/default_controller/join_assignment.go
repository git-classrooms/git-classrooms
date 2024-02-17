package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"log"
)

func (ctrl *DefaultController) JoinAssignment(c *fiber.Ctx) error {
	userID, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomId, err := uuid.Parse(c.Params("classroomId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignmentId, err := uuid.Parse(c.Params("assignmentId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if classroom is valid
	queryClassroom := query.Classroom
	classroom, err := queryClassroom.
		WithContext(c.Context()).
		Where(queryClassroom.ID.Eq(classroomId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if assignment is valid
	queryAssignment := query.Assignment
	assignment, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ID.Eq(assignmentId)).
		Where(queryAssignment.ClassroomID.Eq(classroomId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if invitation is valid
	queryAssignmentProjects := query.AssignmentProjects
	assignmentProject, err := queryAssignmentProjects.
		WithContext(c.Context()).
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

	project, err := repo.ForkProject(assignment.TemplateProjectID, gitlabModel.Private, classroom.GroupID, assignmentName, assignment.Description)
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
