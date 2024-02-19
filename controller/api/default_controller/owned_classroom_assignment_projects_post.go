package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"log"
)

func (ctrl *DefaultController) InviteToAssignmentProject(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	assignment := ctx.GetOwnedClassroomAssignment()

	queryAssignmentProject := query.AssignmentProjects
	assignmentProjects, err := queryAssignmentProject.
		WithContext(c.Context()).
		Where(queryAssignmentProject.AssignmentID.Eq(assignment.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ids := make([]int, len(assignmentProjects))
	for i, project := range assignmentProjects {
		ids[i] = project.UserID
	}

	queryUserClassrooms := query.UserClassrooms
	invitableUsers, err := queryUserClassrooms.WithContext(c.Context()).Preload(queryUserClassrooms.User).
		Where(queryUserClassrooms.ClassroomID.Eq(assignment.ClassroomID)).
		Where(queryUserClassrooms.UserID.NotIn(ids...)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		for _, member := range invitableUsers {
			assignmentProject := &database.AssignmentProjects{
				AssignmentID:       assignment.ID,
				UserID:             member.UserID,
				AssignmentAccepted: false,
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

	repo := context.Get(c).GetGitlabRepository()
	owner, err := repo.GetUserById(classroom.OwnerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, member := range invitableUsers {
		log.Println("Sending invitation to", member.User.GitlabEmail)

		joinPath := fmt.Sprintf("/classrooms/%s/assignments/%s/accept", classroom.ID.String(), assignment.ID.String())
		err = ctrl.mailRepo.SendAssignmentNotification(member.User.GitlabEmail,
			fmt.Sprintf(`You were invited to a new Assigment "%s"`,
				classroom.Name),
			mailRepo.AssignmentNotificationData{
				ClassroomName:      classroom.Name,
				ClassroomOwnerName: owner.Name,
				RecipientName:      member.User.Name,
				AssignmentName:     assignment.Name,
				JoinPath:           joinPath,
			})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}
