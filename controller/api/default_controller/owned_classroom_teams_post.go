package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createOwnedTeamRequest struct {
	Name string `json:"name"`
} //@Name CreateOwnedTeamRequest

func (r createOwnedTeamRequest) isValid() bool {
	return r.Name != ""
}

// @Summary		Create new Team
// @Description	Create a new Team for the given classroom for users to join
// @Id				CreateOwnedClassroomTeam
// @Tags			team
// @Accept			json
// @Param			classroomId		path	string										true	"Classroom ID"	Format(uuid)
// @Param			team			body	default_controller.createOwnedTeamRequest	true	"Classroom Info"
// @Param			X-Csrf-Token	header	string										true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/owned/{classroomId}/teams/{teamId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/owned/{classroomId}/teams [post]
func (ctrl *DefaultController) CreateOwnedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	repo := ctx.GetGitlabRepository()

	requestBody := &createOwnedTeamRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	queryTeam := query.Team
	teams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Where(queryTeam.ClassroomID.Eq(classroom.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if classroom.MaxTeams > 0 && len(teams) >= classroom.MaxTeams {
		return fiber.NewError(fiber.StatusForbidden, "The maximum number of teams has been reached.")
	}

	// reauthenticate the repo with the group access token
	err = repo.GroupAccessLogin(classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	group, err := repo.CreateSubGroup(
		requestBody.Name,
		requestBody.Name,
		classroom.GroupID,
		model.Private,
		fmt.Sprintf("Team %s of classroom %s", requestBody.Name, classroom.Name),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newTeam := &database.Team{
		Name:        requestBody.Name,
		GroupID:     group.ID,
		ClassroomID: classroom.ID,
	}

	err = queryTeam.WithContext(c.Context()).Create(newTeam)
	if err != nil {
		newErr := repo.DeleteGroup(group.ID)
		if newErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s", classroom.ID, newTeam.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
