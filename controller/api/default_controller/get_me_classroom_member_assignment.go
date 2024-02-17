package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"strconv"
)

type getMeClassroomMemberAssignmentResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetMeClassroomMemberAssignment(c *fiber.Ctx) error {
	classroom := context.Get(c).GetClassroom()

	if classroom.Role != database.Owner {
		return fiber.NewError(fiber.StatusForbidden, "only the owner can access the assignments")
	}

	memberId, err := strconv.ParseInt(c.Params("memberId"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no member specified")
	}

	assignmentId, err := uuid.Parse(c.Params("assignmentId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no assignment specified")
	}

	queryAssignment := query.Assignment
	assignment, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ClassroomID.
			Eq(classroom.ClassroomID)).
		Where(queryAssignment.ID.Eq(assignmentId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryAssignmentProjects := query.AssignmentProjects
	assignmentProject, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Where(queryAssignmentProjects.AssignmentID.Eq(assignment.ID)).
		Where(queryAssignmentProjects.UserID.Eq(int(memberId))).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	webURL := ""
	if assignmentProject.AssignmentAccepted {
		repo := context.Get(c).GetGitlabRepository()
		projectFromGitLab, err := repo.GetProjectById(assignmentProject.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		webURL = projectFromGitLab.WebUrl
	}

	response := &getMeClassroomMemberAssignmentResponse{
		AssignmentProjects: *assignmentProject,
		ProjectPath:        webURL,
	}

	return c.JSON(response)
}
