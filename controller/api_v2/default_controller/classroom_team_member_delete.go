package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Remove Member from the team
// @Description	Remove Member from the team
// @Id				RemoveMemberFromTeamV2
// @Tags			member
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			teamId			path	string	true	"Team ID"		Format(uuid)
// @Param			memberId		path	int		true	"Member ID"
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId}/members/{memberId} [delete]
func (ctrl *DefaultController) RemoveMemberFromTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	member := ctx.GetClassroomMember()
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()
	repo := ctx.GetGitlabRepository()

	if classroom.Classroom.MaxTeamSize == 1 {
		return fiber.NewError(fiber.StatusForbidden, "Teams are disabled for this classroom.")
	}

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		member.TeamID = nil
		err := tx.UserClassrooms.WithContext(c.Context()).Save(member)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err := repo.RemoveUserFromGroup(team.GroupID, member.UserID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
