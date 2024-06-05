package factory

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func Classroom(overwrites ...map[string]any) database.Classroom {
	classroom := database.Classroom{}
	classroom.ID = uuid.UUID{}
	classroom.Name = "Test classroom"
	classroom.OwnerID = 1
	classroom.Description = "Classroom description"
	classroom.GroupID = 1
	classroom.GroupAccessTokenID = 20
	classroom.GroupAccessToken = "token"

	mergeOverwrites(classroom, overwrites...)

	return classroom
}

func UserClassroom(userID int, classroomID uuid.UUID, overwrites ...map[string]any) database.UserClassrooms {
	userClassroom := database.UserClassrooms{}
	userClassroom.UserID = userID
	userClassroom.ClassroomID = classroomID

	mergeOverwrites(userClassroom, overwrites...)

	return userClassroom
}

func Invitation(classroomID uuid.UUID, overwrites ...map[string]any) database.ClassroomInvitation {
	invitation := database.ClassroomInvitation{}
	invitation.ID = uuid.New()
	invitation.ClassroomID = classroomID
	invitation.Email = "test@example.com"
	invitation.ExpiryDate = time.Now().Add(24 * time.Hour)
	invitation.Status = database.ClassroomInvitationPending

	mergeOverwrites(invitation, overwrites...)

	return invitation
}

func User(overwrites ...map[string]any) database.User {
	usr := database.User{}
	usr.ID = 1
	usr.GitlabEmail = "test@example.com"
	usr.Name = "Test user"

	mergeOverwrites(usr, overwrites...)

	return usr
}

func AssignmentProject(assignmentID uuid.UUID, teamID uuid.UUID, overwrites ...map[string]any) database.AssignmentProjects {
	project := database.AssignmentProjects{}
	project.TeamID = teamID
	project.AssignmentID = assignmentID
	project.ProjectID = 1

	mergeOverwrites(project, overwrites...)

	return project
}

func Team(classroomID uuid.UUID, overwrites ...map[string]any) database.Team {
	team := database.Team{}
	team.ID = uuid.UUID{}
	team.ClassroomID = classroomID

	mergeOverwrites(team, overwrites...)

	return team
}

func Assignment(classroomID uuid.UUID, overwrites ...map[string]any) database.Assignment {
	assignment := database.Assignment{}
	assignment.ID = uuid.UUID{}
	assignment.ClassroomID = classroomID
	assignment.TemplateProjectID = 1234
	assignment.Name = "Test Assignment"
	assignment.Description = "Test Assignment Description"

	dueDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local).Truncate(time.Second)

	assignment.DueDate = &dueDate

	mergeOverwrites(assignment, overwrites...)

	return assignment
}

func mergeOverwrites(obj any, overwrites ...map[string]any) {
	for _, o := range overwrites {
		merge(obj, o)
	}
}

func merge(obj any, values map[string]any) {
	st := reflect.ValueOf(obj).Elem()

	for k, v := range values {
		f := st.FieldByName(k)
		v := reflect.ValueOf(v)
		f.Set(v)
	}
}
