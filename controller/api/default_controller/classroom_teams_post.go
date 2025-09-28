package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createTeamRequest struct {
	Name string `json:"name"`
} //@Name CreateTeamRequest

func (r createTeamRequest) isValid() bool {
	return r.Name != ""
}

// @Summary		Create new Team
// @Description	Create a new Team for the given classroom and join it if you are a student
// @Id				CreateTeam
// @Tags			team
// @Accept			json
// @Param			classroomId		path	string					true	"Classroom ID"	Format(uuid)
// @Param			team			body	api.createTeamRequest	true	"Classroom Info"
// @Param			X-Csrf-Token	header	string					true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/{classroomId}/teams/{teamId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/teams [post]
func (ctrl *DefaultController) CreateTeam(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	log := ctx.GetLoggerForHandler("CreateTeam")
	userID := ctx.GetUserID()
	classroom := ctx.GetUserClassroom()
	team := classroom.Team
	repo := ctx.GetGitlabRepository()

	if classroom.Classroom.MaxTeamSize == 1 {
		return fiber.NewError(fiber.StatusForbidden, "Teams are disabled for this classroom.")
	}

	if team != nil && classroom.Role == database.Student {
		return fiber.NewError(fiber.StatusForbidden, "You are already a member of a team.")
	}

	if !classroom.Classroom.CreateTeams && classroom.Role == database.Student {
		return fiber.NewError(fiber.StatusForbidden, "Only the owner/moderator can create teams in this classroom.")
	}

	var requestBody createTeamRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
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
	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken, log); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	group, err := repo.CreateSubGroup(
		requestBody.Name,
		requestBody.Name,
		classroom.Classroom.GroupID,
		model.Private,
		fmt.Sprintf("Team %s of classroom %s", requestBody.Name, classroom.Classroom.Name),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteGroup(group.ID); err != nil {
				log.Error("error while deleting gitlab group", "groupID", group.ID, "error", err)
			}
		}
	}()

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

	if err = queryTeam.WithContext(c.Context()).Create(newTeam); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctrl.createAssignmentProjectsIfNeeded(c.Context(), query.Q, &classroom.Classroom, newTeam.ID)

	if _, err = repo.ChangeGroupDescription(group.ID, ctrl.createTeamGitlabDescription(&classroom.Classroom, newTeam.ID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s/teams/%s", classroom.ClassroomID.String(), newTeam.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}

func (ctrl *DefaultController) createTeamGitlabDescription(classroom *database.Classroom, teamID uuid.UUID) string {
	return fmt.Sprintf("%s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s/teams/%s)__", classroom.Description, ctrl.config.PublicURL, classroom.ID.String(), teamID.String())
}

func (ctrl *DefaultController) createAssignmentProjectsIfNeeded(ctx context.Context, tx *query.Query, classroom *database.Classroom, teamID uuid.UUID) error {
	fmt.Println("Team was created trying to invite them to acceptable assignments")
	queryAssignment := tx.Assignment
	assignments, err := queryAssignment.WithContext(ctx).
		Preload(queryAssignment.Projects).
		Where(queryAssignment.ClassroomID.Eq(classroom.ID)).
		Find()
	if err != nil {
		return err
	}
	for _, a := range assignments {
		if a.AcceptableSince == nil || utils.Some(a.Projects, func(ap *database.AssignmentProjects) bool {
			return ap.TeamID == teamID
		}) {
			continue
		}

		fmt.Printf("assignment %s is acceptable inviting team %s\n", a.ID, teamID)

		assignmentProject := &database.AssignmentProjects{
			AssignmentID:  a.ID,
			TeamID:        teamID,
			ProjectStatus: database.Pending,
		}
		if err = tx.AssignmentProjects.WithContext(ctx).Create(assignmentProject); err != nil {
			return err
		}
	}
	return nil
}
