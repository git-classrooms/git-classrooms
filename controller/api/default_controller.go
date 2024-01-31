package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"net/http"
	"net/mail"
	"time"
)

type DefaultController struct {
	mailRepo mailRepo.Repository
}

func NewApiController(mailRepo mailRepo.Repository) *DefaultController {
	return &DefaultController{mailRepo: mailRepo}
}

type CreateClassroomRequest struct {
	Name         string   `json:"name"`
	MemberEmails []string `json:"memberEmails"`
	Description  string   `json:"description"`
}

type InviteToClassroomRequest struct {
	MemberEmails []string `json:"memberEmails"`
}

type CreateAssignmentRequest struct {
	AssigneeUserIds   []int `json:"assigneeUserIds"`
	TemplateProjectId int   `json:"templateProjectId"`
}

func (handler *DefaultController) CreateClassroom(c *fiber.Ctx) error { //TODO: Rework, group-access-token and persisting data is needed
	repo := context.GetGitlabRepository(c)

	var err error
	requestBody := new(CreateClassroomRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	group, err := repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
		requestBody.MemberEmails, //TODO: User shouldn't be added in advance. They should be invited and have the chance not to accept
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, memberEmail := range requestBody.MemberEmails {
		err = repo.CreateGroupInvite(group.ID, memberEmail)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(http.StatusCreated)
}

func (handler *DefaultController) JoinClassroom(c *fiber.Ctx) error {
	invitationIdParameter := c.Params("invitationId")

	invitationId, err := uuid.Parse(invitationIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryClassroomInvitation := query.ClassroomInvitation
	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Where(queryClassroomInvitation.ClassroomID.Eq(invitationId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if time.Now().After(invitation.ExpiryDate) {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	currentUser, err := context.GetGitlabRepository(c).GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomIdParameter := c.Params("classroomId")
	classroomId, err := uuid.Parse(classroomIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if classroomId != invitation.ClassroomID {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if invitation.Email == currentUser.Email {
		repo := context.GetGitlabRepository(c)
		err := repo.AddUserToGroup(invitation.Classroom.GroupID, currentUser.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		queryClassroom := query.UserClassrooms
		err = queryClassroom.Create(&database.UserClassrooms{
			UserID:    currentUser.ID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		})
		if err != nil {
			err := repo.RemoveUserFromGroup(invitation.Classroom.GroupID, currentUser.ID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		_, err = queryClassroomInvitation.WithContext(c.Context()).
			Where(queryClassroomInvitation.ID.Eq(invitation.ID)).
			Update(queryClassroomInvitation.Status, database.InvitationAccepted)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(http.StatusAccepted)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
}

func (handler *DefaultController) InviteToClassroom(c *fiber.Ctx) error {
	user, err := session.Get(c).GetUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryClassroom := query.Classroom
	classroom, err := queryClassroom.
		WithContext(c.Context()).
		Where(queryClassroom.ID.Eq(classroomId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Check if owner or moderator of specific classroom
	if classroom.OwnerID != user.ID {
		queryUserClassroom := query.UserClassrooms
		_, err := queryUserClassroom.
			WithContext(c.Context()).
			Where(queryUserClassroom.ClassroomID.Eq(classroomId)).
			Where(queryUserClassroom.UserID.Eq(user.ID)).
			Where(queryUserClassroom.Role.Eq(uint8(database.Moderator))).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
	}

	requestBody := new(InviteToClassroomRequest)
	if err = c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate email addresses
	validatedEmailAddresses := make([]*mail.Address, len(requestBody.MemberEmails))
	for i, email := range requestBody.MemberEmails {
		validatedEmailAddresses[i], err = mail.ParseAddress(email)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	// Create invitations
	err = query.Q.Transaction(func(tx *query.Query) error {
		for _, email := range validatedEmailAddresses {
			newInvitation := database.ClassroomInvitation{
				Status:      database.InvitationPending,
				ClassroomID: classroomId,
				Email:       email.Address,
				Enabled:     true,
				ExpiryDate:  time.Now().Add(time.Hour * 24 * 14), // Two Weeks, TODO: Add to configuration
			}
			err := tx.ClassroomInvitation.Create(&newInvitation)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send invitations
	for _, email := range validatedEmailAddresses {
		err = handler.mailRepo.Send(email.Address, fmt.Sprintf("Test: New Invitation from Classroom %s", classroom.Name), mailRepo.MailData{}) // TODO: Add meaningful content
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(http.StatusAccepted)
}

func (handler *DefaultController) CreateAssignment(c *fiber.Ctx) error {
	repo := context.GetGitlabRepository(c)

	var err error
	requestBody := new(CreateAssignmentRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	templateProject, err := repo.GetProjectById(requestBody.TemplateProjectId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	name := templateProject.Name

	assignees := make([]model.User, len(requestBody.AssigneeUserIds))
	for i, id := range requestBody.AssigneeUserIds {
		user, err := repo.GetUserById(id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		assignees[i] = *user
		name += "_" + user.Username
	}

	project := &model.Project{}
	project, err = repo.ForkProject(requestBody.TemplateProjectId, name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	project, err = repo.AddProjectMembers(project.ID, assignees)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusCreated)
}
