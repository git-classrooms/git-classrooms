package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) JoinJoinedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	classroom := ctx.GetJoinedClassroom()
	ownTeam := ctx.GetJoinedTeam()
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

	if err := repo.AddUserToGroup(team.GroupID, userID, model.DeveloperPermissions); err != nil {
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
