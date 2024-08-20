package factory

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func Classroom(ownerID int) *database.Classroom {
	classroom := database.Classroom{}
	classroom.Name = gofakeit.Name()
	classroom.OwnerID = ownerID
	classroom.Description = gofakeit.ProductDescription()
	classroom.GroupID = 1
	classroom.GroupAccessTokenID = 20
	classroom.GroupAccessToken = "token"

	err := query.Classroom.WithContext(context.Background()).Create(&classroom)
	if err != nil {
		log.Fatalf("could not insert classroom: %s", err.Error())
	}

	return &classroom
}

func UserClassroom(userID int, classroomID uuid.UUID, role database.Role) *database.UserClassrooms {
	userClassroom := database.UserClassrooms{}
	userClassroom.UserID = userID
	userClassroom.ClassroomID = classroomID
	userClassroom.Role = role

	err := query.UserClassrooms.WithContext(context.Background()).Create(&userClassroom)
	if err != nil {
		log.Fatalf("could not insert classroom: %s", err.Error())
	}

	return &userClassroom
}

func Invitation(classroomID uuid.UUID) *database.ClassroomInvitation {
	invitation := database.ClassroomInvitation{}
	invitation.ClassroomID = classroomID
	invitation.Email = gofakeit.Email()
	invitation.ExpiryDate = time.Now().Add(24 * time.Hour)
	invitation.Status = database.ClassroomInvitationPending

	err := query.ClassroomInvitation.WithContext(context.Background()).Create(&invitation)
	if err != nil {
		log.Fatalf("could not insert invitation: %s", err.Error())
	}

	return &invitation
}

func User() *database.User {
	usr := database.User{}
	usr.GitlabEmail = gofakeit.Email()
	usr.GitlabUsername = gofakeit.Username()
	usr.Name = gofakeit.Name()

	lastUser, err := query.User.WithContext(context.Background()).Last()
	if err != nil {
		usr.ID = 0
	} else {
		usr.ID = lastUser.ID + 1
	}

	err = query.User.WithContext(context.Background()).Create(&usr)
	if err != nil {
		log.Fatalf("could not insert user: %s", err.Error())
	}

	return &usr
}

func AssignmentProject(assignmentID uuid.UUID, teamID uuid.UUID) *database.AssignmentProjects {
	project := database.AssignmentProjects{}
	project.TeamID = teamID
	project.AssignmentID = assignmentID
	project.ProjectID = 1

	err := query.AssignmentProjects.WithContext(context.Background()).Create(&project)
	if err != nil {
		log.Fatalf("could not insert assignment project: %s", err.Error())
	}

	return &project
}

func Team(classroomID uuid.UUID, member []*database.UserClassrooms) *database.Team {
	team := database.Team{}
	team.ClassroomID = classroomID
	team.Member = member

	err := query.Team.WithContext(context.Background()).Create(&team)
	if err != nil {
		log.Fatalf("could not insert team: %s", err.Error())
	}
	return &team
}

func Assignment(classroomID uuid.UUID) *database.Assignment {
	assignment := database.Assignment{}
	assignment.ClassroomID = classroomID
	assignment.TemplateProjectID = 1234
	assignment.Name = "Test Assignment"
	assignment.Description = "Test Assignment Description"

	dueDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local).Truncate(time.Second)

	assignment.DueDate = &dueDate

	err := query.Assignment.WithContext(context.Background()).Create(&assignment)
	if err != nil {
		log.Fatalf("could not insert assignment: %s", err.Error())
	}
	return &assignment
}
