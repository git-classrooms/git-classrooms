package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetActiveAssignments
// @Description	GetActiveAssignments
// @Id				GetActiveAssignments
// @Tags			assignment
// @Produce		json
// @Success		200	{array}		AssignmentResponse
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/assignments [get]
func (ctrl *DefaultController) GetActiveAssignments(c *fiber.Ctx) (err error) {

	ctx := context.Get(c)
	userID := ctx.GetUserID()

	queryAssignment := query.Assignment
	queryUserClassrooms := query.UserClassrooms
	assignments, err := queryAssignment.
		WithContext(c.Context()).
		Join(queryUserClassrooms, queryAssignment.ClassroomID.EqCol(queryUserClassrooms.ClassroomID)).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryAssignment.
			WithContext(c.Context()).
			Where(queryAssignment.DueDate.IsNull()).
			Or(queryAssignment.DueDate.Lt(time.Now())),
		).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(assignments, func(assignment *database.Assignment) *AssignmentResponse {
		return &AssignmentResponse{
			Assignment: assignment,
		}
	})

	return c.JSON(response)
}
