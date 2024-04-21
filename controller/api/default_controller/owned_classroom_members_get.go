package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomMembersResponse struct {
	UserClassrooms []*database.UserClassrooms `json:"classroomMembers"`
	GitlabWebURL   string                     `json:"gitlabWebUrl"`
}

// @Summary		Get classroom Members
// @Description	Create a new classroom
// @Tags			classroom
// @Produces		json
// @Param			classroomId	path		string	true	"Classroom ID" Format(uuid)
// @Success		200			{array}		database.User
// @Failure		401			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/members [get]
func (ctrl *DefaultController) GetOwnedClassroomMembers(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	response := getOwnedClassroomMembersResponse{
		UserClassrooms: classroom.Member,
		GitlabWebURL:   "/api/v1/user/:userId/gitlab",
	}
	return c.JSON(response)
}
