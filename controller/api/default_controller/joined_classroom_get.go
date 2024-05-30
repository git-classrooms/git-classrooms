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
} //@Name GetJoinedClassroomResponse

// @Summary		GetJoinedClassroom
// @Description	GetJoinedClassroom
// @Id				GetJoinedClassroom
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{object}	default_controller.getJoinedClassroomResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/joined/{classroomId} [get]
func (ctrl *DefaultController) GetJoinedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()

	response := &getJoinedClassroomResponse{
		UserClassrooms: *classroom,
		GitlabURL:      fmt.Sprintf("/api/v1/classrooms/joined/%s/gitlab", classroom.ClassroomID.String()),
	}

	return c.JSON(response)
}
