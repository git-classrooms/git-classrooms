package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type filter string //@Name Filter

const (
	ownedClassrooms     filter = "owned"
	moderatorClassrooms filter = "moderator"
	studentClassrooms   filter = "student"
)

type classroomRequestQuery struct {
	Filter filter `query:"filter"`
}

// @Summary		Get classrooms
// @Description	Get classrooms
// @Id				GetClassrooms
// @Tags			classroom
// @Produce		json
// @Param			filter	query		api.filter	false	"Filter Options"
// @Success		200		{array}		api.UserClassroomResponse
// @Failure		401		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Router			/api/v2/classrooms [get]
func (ctrl *DefaultController) GetClassrooms(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	var urlQuery classroomRequestQuery

	if err = c.QueryParser(&urlQuery); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dbQuery := userClassroomQuery(c, userID)
	switch urlQuery.Filter {
	case ownedClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Owner)))
	case moderatorClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Moderator)))
	case studentClassrooms:
		dbQuery = dbQuery.Where(query.UserClassrooms.Role.Eq(uint8(database.Student)))
	default:
	}

	classrooms, err := dbQuery.Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	notArchivedClassrooms := utils.Filter(classrooms, func(classroom *database.UserClassrooms) bool {
		return !classroom.Classroom.Archived || classroom.Role == database.Owner
	})

	response := utils.Map(notArchivedClassrooms, func(classroom *database.UserClassrooms) *UserClassroomResponse {
		return &UserClassroomResponse{
			UserClassrooms: classroom,
			WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/gitlab", classroom.ClassroomID.String()),
		}
	})

	return c.JSON(response)
}
