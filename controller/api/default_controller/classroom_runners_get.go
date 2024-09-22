package api

import (
	"github.com/gofiber/fiber/v2"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type ClassroomRunnerResponse struct {
	*model.Runner
} // @Name ClassroomRunnerResponse

// @Summary		GetClassroomRunners
// @Description	GetClassroomRunners
// @Id				GetClassroomRunners
// @Tags			runners
// @Produce		json
// @Param			classroomId	path		string					true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		ClassroomRunnerResponse	"Success"
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/runners [get]
func (ctrl *DefaultController) GetClassroomRunners(c *fiber.Ctx) (err error) {
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

	allRunners := append(globalRunners, groupRunners...)

	response := make([]ClassroomRunnerResponse, len(allRunners))
	for i, runner := range allRunners {
		response[i] = ClassroomRunnerResponse{Runner: runner}
	}

	return c.JSON(response)
}
