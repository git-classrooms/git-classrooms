package default_controller

import (
	"database/sql/driver"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"strconv"
)

type getMeClassroomMemberAssignmentsResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetMeClassroomMemberAssignments(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetClassroom()

	if classroom.Role != database.Owner {
		return fiber.NewError(fiber.StatusForbidden, "only the owner can access the assignments")
	}

	memberId, err := strconv.ParseInt(c.Params("memberId"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no member specified")
	}

	queryAssignment := query.Assignment
	filteredAssignments, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ClassroomID.
			Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	assignmentIDs := utils.Map(filteredAssignments, func(assignment *database.Assignment) driver.Valuer {
		return assignment.ID
	})

	queryAssignmentProjects := query.AssignmentProjects
	assignmentProjects, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Where(queryAssignmentProjects.AssignmentID.In(assignmentIDs...)).
		Where(queryAssignmentProjects.UserID.Eq(int(memberId))).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	repo := ctx.GetGitlabRepository()
	responses := make([]*getMeClassroomMemberAssignmentsResponse, len(assignmentProjects))
	for i, project := range assignmentProjects {
		webURL := ""
		if project.AssignmentAccepted {
			projectFromGitLab, err := repo.GetProjectById(project.ProjectID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			webURL = projectFromGitLab.WebUrl
		}
		responses[i] = &getMeClassroomMemberAssignmentsResponse{
			AssignmentProjects: *project,
			ProjectPath:        webURL,
		}
	}

	return c.JSON(responses)
}
