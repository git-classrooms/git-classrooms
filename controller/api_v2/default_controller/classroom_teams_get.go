package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeams
// @Description	GetClassroomTeams
// @Id				GetClassroomTeams
// @Tags			team
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		api.TeamResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams [get]
func (ctrl *DefaultController) GetClassroomTeams(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	teams, err := classroomTeamQuery(c, classroom.ClassroomID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(teams, func(team *database.Team) *TeamResponse {
		return &TeamResponse{
			Team:   team,
			WebURL: fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/gitlab", classroom.ClassroomID.String(), team.ID.String()),
		}
	})

	return c.JSON(response)
}
