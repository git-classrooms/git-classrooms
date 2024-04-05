package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetJoinedClassroomTeams(c *fiber.Ctx) error {
	ctx := context.Get(c)

	classroom := ctx.GetJoinedClassroom()

	queryTeam := query.Team
	teams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Where(queryTeam.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(teams)
}
