package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Join the current team
// @Description	Join the current Team if we aren't in another team
// @Id				JoinJoinedClassroomTeam
// @Tags			team
// @Accept			json
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			teamId			path	string	true	"Team ID"		Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/joined/{classroomId}/teams/{teamId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/joined/{classroomId}/teams/{teamId}/join [post]
func (ctrl *DefaultController) JoinJoinedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	classroom := ctx.GetJoinedClassroom()
	ownTeam := classroom.Team
	team := ctx.GetJoinedClassroomTeam()
	repo := ctx.GetGitlabRepository()

	if ownTeam != nil {
		return fiber.NewError(fiber.StatusForbidden, "You are already a member of a team.")
	}

	if len(team.Member) >= classroom.Classroom.MaxTeamSize {
		return fiber.NewError(fiber.StatusForbidden, "The team is full.")
	}

	// reauthenticate the repo with the group access token
	err := repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := repo.AddUserToGroup(team.GroupID, userID, model.ReporterPermissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUserClassrooms := query.UserClassrooms
	_, err = queryUserClassrooms.
		WithContext(c.Context()).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryUserClassrooms.ClassroomID.Eq(classroom.ClassroomID)).
		Update(queryUserClassrooms.TeamID, team.ID)
	if err != nil {
		if err := repo.RemoveUserFromGroup(team.GroupID, userID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/joined/%s/teams/%s", classroom.ClassroomID, team.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
