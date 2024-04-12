package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomTeamResponse struct {
	database.Team
	GitlabUrl string `json:"gitlabUrl"`
}

// @Summary		Get all teams of the current classroom
// @Description	Get all teams of the current classroom
// @Tags			team
// @Accept			json
// @Produces		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		default_controller.getJoinedClassroomTeamResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/joined/{classroomId}/teams [get]
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

	response := utils.Map(teams, func(team *database.Team) *getJoinedClassroomTeamResponse {
		return &getJoinedClassroomTeamResponse{
			Team:      *team,
			GitlabUrl: fmt.Sprintf("/api/v1/classrooms/joined/%s/teams/%s/gitlab", team.ClassroomID.String(), team.ID.String()),
		}
	})

	return c.JSON(response)
}
