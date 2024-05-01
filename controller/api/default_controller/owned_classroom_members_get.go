package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get classroom Members
// @Description	Get classroom Members
// @Id				GetOwnedClassroomMembers
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		default_controller.getOwnedClassroomMemberResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/members [get]
func (ctrl *DefaultController) GetOwnedClassroomMembers(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	members, err := ownedClassroomMemberQuery(classroom.ID, c).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	response := utils.Map(members, func(member *database.UserClassrooms) *getOwnedClassroomMemberResponse {
		return &getOwnedClassroomMemberResponse{
			UserClassrooms: member,
			GitlabURL:      fmt.Sprintf("/api/v1/classrooms/owned/%s/members/%d/gitlab", classroom.ID.String(), member.UserID),
		}
	})
	return c.JSON(response)
}
