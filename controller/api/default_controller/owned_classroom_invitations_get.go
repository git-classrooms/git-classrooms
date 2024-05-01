package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetOwnedClassroomInvitations
// @Description	GetOwnedClassroomInvitations
// @Id				GetOwnedClassroomInvitations
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		ClassroomInvitation
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/classrooms/owned/{classroomId}/invitations [get]
func (ctrl *DefaultController) GetOwnedClassroomInvitations(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	return ctx.JSON(classroom.Invitations)
}
