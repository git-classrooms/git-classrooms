package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomResponse struct {
	database.Classroom
	GitlabURL string `json:"gitlabUrl"`
}

// @Summary		Get classroom
// @Description	Get classroom
// @Id				GetOwnedClassroom
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{object}	default_controller.getOwnedClassroomResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId} [get]
func (ctrl *DefaultController) GetOwnedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	response := &getOwnedClassroomResponse{
		Classroom: *classroom,
		GitlabURL: fmt.Sprintf("/api/v1/classrooms/owned/%s/gitlab", classroom.ID.String()),
	}

	return c.JSON(response)
}
