package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomResponse struct {
	database.UserClassrooms
	GitlabURL string `json:"gitlabUrl"`
}

// @Summary		GetJoinedClassroom
// @Description	GetJoinedClassroom
// @Id				GetJoinedClassroom
// @Tags			assignment
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{object}	default_controller.getJoinedClassroomResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/joined/{classroomId} [get]
func (ctrl *DefaultController) GetJoinedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()

	response := &getJoinedClassroomResponse{
		UserClassrooms: *classroom,
		GitlabURL:      fmt.Sprintf("/api/v1/classrooms/joined/%s/gitlab", classroom.ClassroomID.String()),
	}

	return c.JSON(response)
}
