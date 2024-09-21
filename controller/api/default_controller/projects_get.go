package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type activeAssignmentResponse struct {
	*database.Assignment
	Classroom *database.Classroom `json:"classroom"`
} //	@Name	ActiveAssignmentResponse

type activeAssignmentRequestQuery struct {
	Filter filter `query:"filter"`
}

//	@Summary		GetActiveAssignments
//	@Description	GetActiveAssignments
//	@Id				GetActiveAssignments
//	@Tags			assignment
//	@Produce		json
//	@Param			filter	query		api.filter	false	"Filter Options"
//	@Success		200		{array}		ActiveAssignmentResponse
//	@Failure		401		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/assignments [get]
func (ctrl *DefaultController) GetActiveAssignments(c *fiber.Ctx) (err error) {

	ctx := context.Get(c)
	userID := ctx.GetUserID()

	var urlQuery activeAssignmentRequestQuery

	if err = c.QueryParser(&urlQuery); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryAssignment := query.Assignment
	queryUserClassrooms := query.UserClassrooms

	dbQuery := queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Classroom).
		Join(queryUserClassrooms, queryAssignment.ClassroomID.EqCol(queryUserClassrooms.ClassroomID)).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryAssignment.
			WithContext(c.Context()).
			Where(queryAssignment.DueDate.IsNull()).
			Or(queryAssignment.DueDate.Gt(time.Now())),
		)

	switch urlQuery.Filter {
	case ownedClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Owner)))
	case moderatorClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Moderator)))
	case studentClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Student)))
	default:
	}

	assignments, err := dbQuery.Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(assignments, func(assignment *database.Assignment) *activeAssignmentResponse {
		return &activeAssignmentResponse{
			Assignment: assignment,
			Classroom:  &assignment.Classroom,
		}
	})

	return c.JSON(response)
}
