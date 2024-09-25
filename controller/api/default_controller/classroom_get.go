package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get classroom
// @Description	Get classroom
// @Id				GetClassroom
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{object}	api.UserClassroomResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId} [get]
func (ctrl *DefaultController) GetClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	response := &UserClassroomResponse{
		UserClassrooms:   classroom,
		WebURL:           fmt.Sprintf("/api/v1/classrooms/%s/gitlab", classroom.ClassroomID.String()),
		AssignmentsCount: len(classroom.Classroom.Assignments),
	}

	return c.JSON(response)
}
