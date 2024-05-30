package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomMembers
// @Description	GetClassroomMembers
// @Id				GetClassroomMembers
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		api.UserClassroomResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/members [get]
func (ctrl *DefaultController) GetClassroomMembers(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	members, err := classroomMemberQuery(c, classroom.ClassroomID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(members, func(member *database.UserClassrooms) *UserClassroomResponse {
		return &UserClassroomResponse{
			UserClassrooms: member,
			WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/members/%d/gitlab", classroom.ClassroomID.String(), member.UserID),
		}
	})

	return c.JSON(response)
}
