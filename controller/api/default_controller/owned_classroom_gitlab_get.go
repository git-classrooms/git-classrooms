package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get all Projects of current team
// @Description	Get all gitlab projects of the current team
// @Tags			classroom
// @Accept			json
// @Param			classroomId	path	string	true	"Classroom ID"	Format(uuid)
// @Success		302
// @Header			302	{string}	Location	"<Gitlab Group url>"
// @Failure		400	{object}	httputil.HTTPError
// @Failure		401	{object}	httputil.HTTPError
// @Failure		404	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/gitlab [get]
func (ctrl *DefaultController) GetOwnedClassroomGitlab(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	repo := ctx.GetGitlabRepository()

	group, err := repo.GetGroupById(classroom.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect(group.WebUrl)
}
