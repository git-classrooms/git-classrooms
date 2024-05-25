package api

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateTeamRequest struct {
	Name string `json:"name,omitempty"`
} //@Name UpdateTeamRequest

func (r updateTeamRequest) isValid() bool {
	return r.Name != ""
}

// @Summary		Update Team
// @Description	Update Team
// @Id				UpdateTeam
// @Tags			team
// @Produces		json
// @Param			classroomId		path	string					true	"Classroom ID"	Format(uuid)
// @Param			teamId			path	string					true	"Team ID"		Format(uuid)
// @Param			UpdateTeam		body	api.updateTeamRequest	true	"Update Team"
// @Param			X-Csrf-Token	header	string					true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId} [put]
func (ctrl *DefaultController) UpdateTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()
	repo := ctx.GetGitlabRepository()

	if classroom.Classroom.MaxTeamSize == 1 {
		return fiber.NewError(fiber.StatusForbidden, "Teams are disabled for this classroom.")
	}

	if classroom.Role == database.Student && team.ClassroomID != classroom.ClassroomID {
		return fiber.NewError(fiber.StatusForbidden, "You are not a member of this team.")
	}

	var requestBody updateTeamRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	oldTeamName := team.Name

	_, err = repo.ChangeGroupName(team.GroupID, requestBody.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		if recover() != nil || err != nil {
			if _, err := repo.ChangeGroupName(team.GroupID, oldTeamName); err != nil {
				log.Printf("Failed to revert group name: %v", err)
			}
		}
	}()

	team.Name = requestBody.Name

	err = query.Team.WithContext(c.Context()).Save(team)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusAccepted)
}
