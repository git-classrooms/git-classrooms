package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createJoinedTeamRequest struct {
	Name string `json:"name"`
} //@Name CreateJoinedTeamRequest

func (r createJoinedTeamRequest) isValid() bool {
	return r.Name != ""
}

// @Summary		Create new Team
// @Description	Create a new Team for the given classroom and join it if you are not a moderator
// @Id				CreateJoinedClassroomTeam
// @Tags			team
// @Accept			json
// @Param			classroomId		path	string										true	"Classroom ID"	Format(uuid)
// @Param			team			body	default_controller.createJoinedTeamRequest	true	"Classroom Info"
// @Param			X-Csrf-Token	header	string										true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/joined/{classroomId}/teams/{teamId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/joined/{classroomId}/teams [post]
func (ctrl *DefaultController) CreateJoinedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	classroom := ctx.GetJoinedClassroom()
	team := classroom.Team
	repo := ctx.GetGitlabRepository()

	if team != nil && classroom.Role != database.Moderator {
		return fiber.NewError(fiber.StatusForbidden, "You are already a member of a team.")
	}

	requestBody := &createJoinedTeamRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if !classroom.Classroom.CreateTeams && classroom.Role != database.Moderator {
		return fiber.NewError(fiber.StatusForbidden, "Only the owner/moderator can create teams in this classroom.")
	}

	queryTeam := query.Team
	teams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Where(queryTeam.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if classroom.Classroom.MaxTeams > 0 && len(teams) >= classroom.Classroom.MaxTeams {
		return fiber.NewError(fiber.StatusForbidden, "The maximum number of teams has been reached.")
	}

	// reauthenticate the repo with the group access token
	err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	group, err := repo.CreateSubGroup(
		requestBody.Name,
		classroom.Classroom.GroupID,
		model.Private,
		fmt.Sprintf("Team %s of classroom %s", requestBody.Name, classroom.Classroom.Name),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUserClassrooms := query.UserClassrooms
	user, err := queryUserClassrooms.
		WithContext(c.Context()).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryUserClassrooms.ClassroomID.Eq(classroom.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	member := make([]*database.UserClassrooms, 0)

	if classroom.Role == database.Student {
		if err := repo.AddUserToGroup(group.ID, userID, model.ReporterPermissions); err != nil {
			if err := repo.DeleteGroup(group.ID); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		member = append(member, user)
	}

	newTeam := &database.Team{
		Name:        requestBody.Name,
		GroupID:     group.ID,
		ClassroomID: classroom.ClassroomID,
		Member:      member,
	}

	err = queryTeam.WithContext(c.Context()).Create(newTeam)
	if err != nil {
		newErr := repo.DeleteGroup(group.ID)
		if newErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/joined/%s/teams/%s", classroom.ClassroomID, newTeam.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
