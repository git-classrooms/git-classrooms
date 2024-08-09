package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type ActiveProjectResponse struct {
	*database.AssignmentProjects
	WebURL      string `json:"webUrl"`
	ClassroomID string `json:"classroomId"`
} //@Name ActiveProjectResponse

// @Summary		GetActiveProjects
// @Description	GetActiveProjects
// @Id				GetActiveProjects
// @Tags			project
// @Produce		json
// @Success		200	{array}		ActiveProjectResponse
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/projects [get]
func (ctrl *DefaultController) GetActiveProjects(c *fiber.Ctx) (err error) {

	ctx := context.Get(c)
	userID := ctx.GetUserID()

	queryAssignmentProjects := query.AssignmentProjects
	queryAssignment := query.Assignment
	queryTeam := query.Team
	queryUserClassrooms := query.UserClassrooms
	projects, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Preload(queryAssignmentProjects.Assignment).
		Join(queryTeam, queryAssignmentProjects.TeamID.EqCol(queryTeam.ID)).
		Join(queryAssignment, queryAssignmentProjects.AssignmentID.EqCol(queryAssignment.ID)).
		Join(queryUserClassrooms, queryTeam.ID.EqCol(queryUserClassrooms.TeamID)).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryAssignmentProjects.
			WithContext(c.Context()).
			Where(queryAssignment.DueDate.IsNull()).
			Or(queryAssignment.DueDate.Lt(time.Now())),
		).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(projects, func(project *database.AssignmentProjects) *ActiveProjectResponse {
		return &ActiveProjectResponse{
			AssignmentProjects: project,
			WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/projects/%s/gitlab", project.Assignment.ClassroomID.String(), project.ID.String()),
			ClassroomID:        project.Assignment.ClassroomID.String(),
		}
	})

	return c.JSON(response)
}
