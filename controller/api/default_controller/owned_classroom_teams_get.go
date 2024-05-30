package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get all teams of the current classroom
// @Description	Get all teams of the current classroom
// @Id				GetOwnedClassroomTeams
// @Tags			team
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		default_controller.getOwnedClassroomTeamResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/owned/{classroomId}/teams [get]
func (ctrl *DefaultController) GetOwnedClassroomTeams(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	teams, err := ownedClassroomTeamQuery(c, classroom.ID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(teams, func(team *database.Team) *getOwnedClassroomTeamResponse {
		member := utils.Map(team.Member, func(u *database.UserClassrooms) *database.User {
			return &u.User
		})
		return &getOwnedClassroomTeamResponse{
			Team:       *team,
			UserMember: member,
			GitlabURL:  fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/gitlab", classroom.ID.String(), team.ID.String()),
		}
	})

	return c.JSON(response)
}
