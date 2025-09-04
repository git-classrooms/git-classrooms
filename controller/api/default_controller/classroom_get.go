package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
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

	var inviteCode *uuid.UUID
	if classroom.Role <= database.Moderator {
		inviteCode = &classroom.Classroom.InviteCode
	}

	response := &UserClassroomResponse{
		UserClassrooms:   classroom,
		WebURL:           fmt.Sprintf("/api/v1/classrooms/%s/gitlab", classroom.ClassroomID.String()),
		AssignmentsCount: len(classroom.Classroom.Assignments),
		InviteCode:       inviteCode,
	}

	return c.JSON(response)
}
