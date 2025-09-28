package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
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
	repo := ctx.GetGitlabRepository()

	if classroom.Classroom.MaxTeamSize == 1 {
		return fiber.NewError(fiber.StatusForbidden, "Teams are disabled for this classroom.")
	}

	if classroom.Team != nil && classroom.Role == database.Student {
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
		return fiber.NewError(fiber.StatusBadRequest, "Body is not valid")
	}

	log.Debug("retrieving all teams from db")

	queryTeam := query.Team
	teams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Where(queryTeam.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		log.Error("error while retrieving teams from db", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "database error")
	}

	if classroom.Classroom.MaxTeams > 0 && len(teams) >= classroom.Classroom.MaxTeams {
		return fiber.NewError(fiber.StatusForbidden, "The maximum number of teams has been reached.")
	}

	log.Debug("authenticate the repo with group access token")

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
		log.Error("error creating gitlab subgroup", "groupID", classroom.Classroom.GroupID, "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteGroup(group.ID); err != nil {
				log.Error("error while deleting gitlab group", "groupID", group.ID, "error", err)
			}
		}
	}()

	log.Info("created subgroup in gitlab for team", "parentGroup", classroom.Classroom.GroupID, "subGroupID", group.ID, "name", requestBody.Name)

	member := make([]*database.UserClassrooms, 0)

	if classroom.Role == database.Student {
		log.Info("adding user to gitlab group because you are a student", "groupID", group.ID)
		if err = repo.AddUserToGroup(group.ID, userID, model.ReporterPermissions); err != nil {
			log.Error("error adding user to group", "groupID", group.ID, "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		member = append(member, classroom)
	}

	newTeam := &database.Team{
		Name:        requestBody.Name,
		GroupID:     group.ID,
		ClassroomID: classroom.ClassroomID,
		Member:      member,
	}

	if err = queryTeam.WithContext(c.Context()).Create(newTeam); err != nil {
		log.Error("error creating team in db", "name", requestBody.Name, "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log = log.With("team", newTeam)
	ctx.SetLogger(log)

	ctrl.createAssignmentProjectsIfNeeded(c.Context(), query.Q, &classroom.Classroom, newTeam.ID)

	log.Debug("changing description of gitlab group", "groupID", group.ID)

	if _, err = repo.ChangeGroupDescription(group.ID, ctrl.createTeamGitlabDescription(&classroom.Classroom, newTeam.ID)); err != nil {
		log.Error("error changing group description", "groupID", group.ID, "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s/teams/%s", classroom.ClassroomID.String(), newTeam.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}

func (ctrl *DefaultController) createTeamGitlabDescription(classroom *database.Classroom, teamID uuid.UUID) string {
	return fmt.Sprintf("%s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s/teams/%s)__", classroom.Description, ctrl.config.PublicURL, classroom.ID.String(), teamID.String())
}

func (ctrl *DefaultController) createAssignmentProjectsIfNeeded(ctx context.Context, tx *query.Query, classroom *database.Classroom, teamID uuid.UUID) error {
	log := logging.GetLogger(ctx)
	log.Debug("team was created trying to invite them to acceptable assignments")

	queryAssignment := tx.Assignment
	assignments, err := queryAssignment.WithContext(ctx).
		Preload(queryAssignment.Projects).
		Where(queryAssignment.ClassroomID.Eq(classroom.ID)).
		Find()
	if err != nil {
		log.Error("error getting assignments from db", "error", err)
		return err
	}
	for _, a := range assignments {
		if a.AcceptableSince == nil || utils.Some(a.Projects, func(ap *database.AssignmentProjects) bool {
			return ap.TeamID == teamID
		}) {
			continue
		}

		log.Info("assignment is acceptable inviting team", "assignment", a)

		assignmentProject := &database.AssignmentProjects{
			AssignmentID:  a.ID,
			TeamID:        teamID,
			ProjectStatus: database.Pending,
		}
		if innerErr := tx.AssignmentProjects.WithContext(ctx).Create(assignmentProject); innerErr != nil {
			log.Error("error creating assignmentProject", "error", innerErr)
			err = innerErr
			continue
		}
	}
	return err
}
