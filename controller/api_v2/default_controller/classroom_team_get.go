package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	response := &TeamResponse{
		Team:   team,
		WebURL: fmt.Sprintf("/api/v2/classrooms/%s/teams/%s", classroom.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
