package default_controller

import (
	"fmt"
	"gorm.io/gen/field"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		InviteToAssignment
// @Description	InviteToAssignment
// @Id				InviteToAssignment
// @Tags			assignment
// @Accept			json
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string	true	"Assignment ID"	Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		201
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/owned/{classroomId}/assignments/{assignmentId}/projects [post]
func (ctrl *DefaultController) InviteToAssignmentProject(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	classroom := ctx.GetOwnedClassroom()
	assignment := ctx.GetOwnedClassroomAssignment()

	queryAssignmentProject := query.AssignmentProjects
	assignmentProjects, err := queryAssignmentProject.
		WithContext(c.Context()).
		Preload(queryAssignmentProject.Team).
		Preload(field.NewRelation("Team.Member", "")).
		Where(queryAssignmentProject.AssignmentID.Eq(assignment.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ids := make([]uuid.UUID, len(assignmentProjects))
	for i, project := range assignmentProjects {
		ids[i] = project.TeamID
	}

	queryTeam := query.Team
	invitableTeams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Preload(queryTeam.Member.User).
		FindByClassroomIDAndNotInTeamIDs(classroom.ID, ids)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		for _, team := range invitableTeams {
			assignmentProject := &database.AssignmentProjects{
				AssignmentID:  assignment.ID,
				TeamID:        team.ID,
				ProjectStatus: database.Pending,
			}
			if err := tx.AssignmentProjects.WithContext(c.Context()).Create(assignmentProject); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userQuery := query.User
	me, err := userQuery.
		WithContext(c.Context()).
		Where(userQuery.ID.Eq(userID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, team := range invitableTeams {
		for _, member := range team.Member {
			log.Println("Sending invitation to", member.User.GitlabEmail)

			joinPath := fmt.Sprintf("/classrooms/joined/%s/assignments/%s/accept", classroom.ID.String(), assignment.ID.String())
			err = ctrl.mailRepo.SendAssignmentNotification(member.User.GitlabEmail,
				fmt.Sprintf(`You were invited to a new Assigment "%s"`,
					classroom.Name),
				mailRepo.AssignmentNotificationData{
					ClassroomName:      classroom.Name,
					ClassroomOwnerName: me.Name,
					RecipientName:      member.User.Name,
					AssignmentName:     assignment.Name,
					JoinPath:           joinPath,
				})
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}
