package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get your classrooms
// @Description	Get your classrooms
// @Id				GetClassrooms
// @Tags			classroom
// @Produce		json
// @Success		200	{array}		api.UserClassroomResponse
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms [get]
func (ctrl *DefaultController) GetClassrooms(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	classrooms, err := userClassroomQuery(c, userID).
		Find()
	if err != nil {
		return err
	}

	response := utils.Map(classrooms, func(classroom *database.UserClassrooms) *UserClassroomResponse {
		return &UserClassroomResponse{
			UserClassrooms: classroom,
			WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/gitlab", classroom.ClassroomID.String()),
		}
	})

	return c.JSON(response)
}
