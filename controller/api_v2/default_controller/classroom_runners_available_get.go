package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomRunnersAreAvailable
// @Description	GetClassroomRunnersAreAvailable
// @Id				GetClassroomRunnersAreAvailable
// @Tags			runners
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200		{boolean}	true	"Success"
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/runners/available [get]
func (ctrl *DefaultController) GetClassroomRunnersAreAvailable(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repository := ctx.GetGitlabRepository()
	classroom := ctx.GetUserClassroom()

	globalRunners, err := repository.GetAvailableRunnersForGitLab()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	groupRunners, err := repository.GetAvailableRunnersForGroup(classroom.Classroom.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(globalRunners) == 0 && len(groupRunners) == 0 {
		return c.JSON(false)
	} else {
		return c.JSON(true)
	}
}
