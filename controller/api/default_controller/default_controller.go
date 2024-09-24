package api

import (
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
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
	g        *singleflight.Group
}

func NewAPIV1Controller(mailRepo mailRepo.Repository, config config.ApplicationConfig) *DefaultController {
	g := &singleflight.Group{}
	return &DefaultController{mailRepo: mailRepo, config: config, g: g}
}

type UserResponse struct {
	*database.User
	WebURL string `json:"webUrl"`
} // @Name UserResponse

type TeamResponse struct {
	*database.Team
	Members []*UserClassroomResponse `json:"members"`
	WebURL  string                   `json:"webUrl"`
} // @Name TeamResponse

type AssignmentResponse struct {
	*database.Assignment
} // @Name AssignmentResponse

type ProjectResponse struct {
	*database.AssignmentProjects
	WebURL       string `json:"webUrl"`
	ReportWebURL string `json:"reportWebUrl"`
} // @Name ProjectResponse

type UserClassroomResponse struct {
	*database.UserClassrooms
	WebURL           string `json:"webUrl"`
	AssignmentsCount int    `json:"assignmentsCount"`
} // @Name UserClassroomResponse
