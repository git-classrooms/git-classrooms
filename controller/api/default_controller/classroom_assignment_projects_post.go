package api

import (
	"database/sql/driver"
	"fmt"
	"log"

	"gorm.io/gen/field"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		InviteToAssignment
// @Description	InviteToAssignment
// @Id				InviteToAssignment
// @Tags			project
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
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/projects [post]
func (ctrl *DefaultController) InviteToAssignment(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	classroom := ctx.GetUserClassroom()
	assignment := ctx.GetAssignment()

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

	ids := utils.Map(assignmentProjects, func(p *database.AssignmentProjects) driver.Valuer {
		return p.TeamID
	})

	queryTeam := query.Team
	invitableTeams, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Preload(queryTeam.Member.User).
		Where(queryTeam.ClassroomID.Eq(classroom.ClassroomID)).
		Not(queryTeam.ID.In(ids...)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) (err error) {
		for _, team := range invitableTeams {
			assignmentProject := &database.AssignmentProjects{
				AssignmentID:  assignment.ID,
				TeamID:        team.ID,
				ProjectStatus: database.Pending,
			}
			if err = tx.AssignmentProjects.WithContext(c.Context()).Create(assignmentProject); err != nil {
				return err
			}
			team.AssignmentProjects = []*database.AssignmentProjects{assignmentProject}
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUser := query.User
	me, err := queryUser.
		WithContext(c.Context()).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, team := range invitableTeams {
		for _, member := range team.Member {
			log.Println("Sending invitation to", member.User.GitlabEmail)

			joinPath := fmt.Sprintf("/classrooms/%s/projects/%s/accept", classroom.ClassroomID.String(), team.AssignmentProjects[0].ID.String())
			err = ctrl.mailRepo.SendAssignmentNotification(member.User.GitlabEmail,
				fmt.Sprintf(`You were invited to a new Assigment "%s"`,
					classroom.Classroom.Name),
				mailRepo.AssignmentNotificationData{
					ClassroomName:      classroom.Classroom.Name,
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
