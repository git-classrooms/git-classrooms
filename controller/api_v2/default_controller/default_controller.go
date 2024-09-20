package api

import (
	"gitlab.hs-flensburg.de/gitlab-classroom/config"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
)

type Params struct {
	ClassroomID         *uuid.UUID `params:"classroomId"`
	AssignmentID        *uuid.UUID `params:"assignmentId"`
	AssignmentProjectID *uuid.UUID `params:"projectId"`
	MemberID            *int       `params:"memberId"`
	TeamID              *uuid.UUID `params:"teamId"`
	InvitationID        *uuid.UUID `params:"invitationId"`
}

type DefaultController struct {
	config   config.ApplicationConfig
	mailRepo mailRepo.Repository
}

func NewApiV2Controller(mailRepo mailRepo.Repository, config config.ApplicationConfig) *DefaultController {
	return &DefaultController{mailRepo: mailRepo, config: config}
}

type UserResponse struct {
	*database.User
	WebURL string `json:"webUrl"`
} //@Name UserResponse

type TeamResponse struct {
	*database.Team
	Members []*UserClassroomResponse `json:"members"`
	WebURL  string                   `json:"webUrl"`
} //@Name TeamResponse

type AssignmentResponse struct {
	*database.Assignment
} //@Name AssignmentResponse

type ProjectResponse struct {
	*database.AssignmentProjects
	WebURL       string `json:"webUrl"`
	ReportWebURL string `json:"reportWebUrl"`
} //@Name ProjectResponse

type UserClassroomResponse struct {
	*database.UserClassrooms
	WebURL           string `json:"webUrl"`
	AssignmentsCount int    `json:"assignmentsCount"`
} //@Name UserClassroomResponse
