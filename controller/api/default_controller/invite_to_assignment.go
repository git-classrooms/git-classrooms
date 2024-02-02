package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"log"
)

// /classrooms/:classroomId/assignments/:assignmentId
func (ctrl *DefaultController) InviteToAssignment(c *fiber.Ctx) error {
	userID, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomID, err := uuid.Parse(c.Params("classroomID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignmentID, err := uuid.Parse(c.Params("assignmentID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if classroom is valid
	queryClassroom := query.Classroom
	classroom, err := queryClassroom.
		WithContext(c.Context()).
		Preload(queryClassroom.Member).
		Where(queryClassroom.ID.Eq(classroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Check if owner or moderator of specific classroom
	if classroom.OwnerID != userID {
		queryUserClassroom := query.UserClassrooms
		_, err := queryUserClassroom.
			WithContext(c.Context()).
			Where(queryUserClassroom.ClassroomID.Eq(classroomID)).
			Where(queryUserClassroom.UserID.Eq(userID)).
			Where(queryUserClassroom.Role.Eq(uint8(database.Moderator))).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
	}

	// Check if assignment is valid
	queryAssignment := query.Assignment
	assigment, err := queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Projects).
		Preload(queryAssignment.Projects.User).
		Where(queryAssignment.ID.Eq(assignmentID)).
		Where(queryAssignment.ClassroomID.Eq(classroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Fetch candidates for invitation mail
	queryUser := query.User
	queryUserClassrooms := query.UserClassrooms
	queryAssignmentProject := query.AssignmentProjects

	invitableUsers, err := queryUser.
		WithContext(c.Context()).
		Join(queryUserClassrooms, queryUser.ID.EqCol(queryUserClassrooms.UserID)).
		LeftJoin(queryAssignmentProject, queryAssignmentProject.UserID.EqCol(queryUser.ID)).
		Where(queryUserClassrooms.ClassroomID.Eq(classroomID)).
		Where(queryAssignmentProject.ID.IsNull()).
		Where(queryUserClassrooms.UserID.Neq(classroom.OwnerID)).
		Find()
	if err != nil {
		return err
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		for _, member := range invitableUsers {
			assignmentProject := &database.AssignmentProjects{
				AssignmentID:       assignmentID,
				UserID:             member.ID,
				AssignmentAccepted: false,
			}
			if err := tx.AssignmentProjects.Create(assignmentProject); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	repo := context.GetGitlabRepository(c)
	owner, err := repo.GetUserById(classroom.OwnerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, member := range invitableUsers {
		log.Println("Sending invitation to", member.GitlabEmail)

		joinPath := fmt.Sprintf("/api/classrooms/%s/assignments/%s/accept", classroom.ID.String(), assigment.ID.String())
		err = ctrl.mailRepo.SendAssignmentNotification(member.GitlabEmail,
			fmt.Sprintf(`Test: You were invited to a new Assigment "%s"`,
				classroom.Name),
			mailRepo.AssignmentNotificationData{
				ClassroomName:      classroom.Name,
				ClassroomOwnerName: owner.Name,
				RecipientName:      member.Name,
				AssignmentName:     assigment.Name,
				JoinPath:           joinPath,
			})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}

func filterInvitableUsers(users []*database.User, assigmentProjects []*database.AssignmentProjects) []*database.User {
	invitableUsers := make([]*database.User, 0)
	for _, user := range users {
		for _, assignmentProject := range assigmentProjects {
			if user.ID != assignmentProject.UserID {
				invitableUsers = append(invitableUsers, user)
			}
		}
	}

	return invitableUsers
}
