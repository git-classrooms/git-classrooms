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
// @Success		200			{array}		database.ClassroomInvitation
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		403			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/invitations [get]
func (ctrl *DefaultController) GetOwnedClassroomInvitations(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	return ctx.JSON(classroom.Invitations)
}
