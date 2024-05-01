package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomMemberResponse struct {
	*database.UserClassrooms
	GitlabURL string `json:"gitlabUrl"`
}

// @Summary		Get classroom Member
// @Description	Get classroom Member
// @Id				GetOwnedClassroomMember
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			memberId	path		int		true	"Member ID"
// @Success		200			{object}	default_controller.getOwnedClassroomMemberResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/members [get]
func (ctrl *DefaultController) GetOwnedClassroomMember(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	member := ctx.GetOwnedClassroomMember()
	response := &getOwnedClassroomMemberResponse{
		UserClassrooms: member,
		GitlabURL:      fmt.Sprintf("/api/v1/classrooms/owned/%s/members/%d/gitlab", classroom.ID.String(), member.UserID),
	}

	return c.JSON(response)
}
