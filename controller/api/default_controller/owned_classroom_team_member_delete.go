package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

//	@Summary		Remove current Member from the current team
//	@Description	Remove current Member from the current team
//	@Tags			team, member
//	@Accept			json
//	@Param			classroomId	path	string	true	"Classroom ID"	Format(uuid)
//	@Param			teamId		path	string	true	"Team ID"		Format(uuid)
//	@Param			memberId	path	string	true	"Member ID"		Format(uuid)
//	@Success		204
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		401	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/classrooms/owned/{classroomId}/teams/{teamdId}/members/{memberId} [delete]
func (ctrl *DefaultController) RemoveMemberFromTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	member := ctx.GetOwnedClassroomTeamMember()
	classroom := ctx.GetOwnedClassroom()
	team := ctx.GetOwnedClassroomTeam()
	repo := ctx.GetGitlabRepository()

	if err := repo.GroupAccessLogin(classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err := query.Q.Transaction(func(tx *query.Query) error {
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
