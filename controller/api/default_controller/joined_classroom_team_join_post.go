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

	// reauthenticate the repo with the group access token
	err := repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := repo.AddUserToGroup(team.GroupID, userID, model.DeveloperPermissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := query.User.WithContext(c.Context()).Where(query.User.ID.Eq(userID)).First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryTeam := query.Team
	err = queryTeam.
		Member.Model(team).
		Append(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/joined/%s/teams/%s", classroom.ClassroomID, team.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
