package api

import (
	"context"
	"fmt"

	"io"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestPutOwnedAssignments(t *testing.T) {
	restoreDatabase(t)

	owner := factory.User()
	student1 := factory.User()
	student2 := factory.User()

	classroom := factory.Classroom(owner.ID)

	dueDate := time.Now().Add(1 * time.Hour)
	dueDate = dueDate.Truncate(time.Second)

	assignment := factory.Assignment(classroom.ID, &dueDate, false)

	members := []*database.UserClassrooms{
		factory.UserClassroom(owner.ID, classroom.ID, database.Owner),
		factory.UserClassroom(student1.ID, classroom.ID, database.Student),
		factory.UserClassroom(student2.ID, classroom.ID, database.Student),
	}

	team := factory.Team(classroom.ID, members)
	project := factory.AssignmentProject(assignment.ID, team.ID)

	project.ProjectStatus = database.Pending
	query.AssignmentProjects.WithContext(context.Background()).Save(project)

	app, gitlabRepo, _ := setupApp(t, owner)
	targetRoute := fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

	t.Run("updates assignment", func(t *testing.T) {
		newTime := time.Now().Add(time.Hour * 24)
		newTime = newTime.Truncate(time.Second)

		requestBody := updateAssignmentRequest{
			Name:        "New",
			Description: "new",
			DueDate:     &newTime,
		}

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		updatedAssignment, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment.ID)).
			First()

		assert.NoError(t, err)
		assert.Equal(t, requestBody.Name, updatedAssignment.Name)
		assert.Equal(t, requestBody.Description, updatedAssignment.Description)
		assert.Equal(t, newTime, *updatedAssignment.DueDate)
	})

	t.Run("due date is in the past", func(t *testing.T) {
		newTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
		requestBody := updateAssignmentRequest{
			DueDate: &newTime,
		}

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		assert.Equal(t, "DueDate must be in the future", bodyString)
	})

	t.Run("assignment name and description can not be changed after it has been accepted by students", func(t *testing.T) {
		newTime := time.Now().Add(time.Hour * 24)
		requestBody := updateAssignmentRequest{
			Name:        "New",
			Description: "new",
			DueDate:     &newTime,
		}

		project.ProjectStatus = database.Accepted
		query.AssignmentProjects.WithContext(context.Background()).Save(project)

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		assert.Equal(t, "Assignment name and description can not be changed after it has been accepted by students", bodyString)
	})

	t.Run("reopenAssignment", func(t *testing.T) {
		newTime := time.Now().Add(time.Hour * 26)
		newTime = newTime.Truncate(time.Second)
		requestBody := updateAssignmentRequest{
			DueDate: &newTime,
		}

		project.ProjectStatus = database.Accepted
		query.AssignmentProjects.WithContext(context.Background()).Save(project)

		assignment.Closed = true
		query.Assignment.WithContext(context.Background()).Save(assignment)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, owner.ID).
			Return(model.OwnerPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, student1.ID).
			Return(model.ReporterPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, student2.ID).
			Return(model.ReporterPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			ChangeUserAccessLevelInProject(project.ProjectID, student1.ID, model.DeveloperPermissions).
			Return(nil).
			Times(1)

		gitlabRepo.EXPECT().
			ChangeUserAccessLevelInProject(project.ProjectID, student2.ID, model.DeveloperPermissions).
			Return(nil).
			Times(1)

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		gitlabRepo.AssertExpectations(t)

		updatedAssignment, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment.ID)).
			First()
		assert.NoError(t, err)
		assert.False(t, updatedAssignment.Closed)
		assert.Equal(t, newTime, *updatedAssignment.DueDate)
	})

	t.Run("repo.ChangeUserAccessLevelInProject throws error", func(t *testing.T) {
		newTime := time.Now().Add(time.Hour * 25)
		newTime = newTime.Truncate(time.Second)
		requestBody := updateAssignmentRequest{
			DueDate: &newTime,
		}

		project.ProjectStatus = database.Accepted
		query.AssignmentProjects.WithContext(context.Background()).Save(project)

		assignment.Closed = true
		query.Assignment.WithContext(context.Background()).Save(assignment)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, owner.ID).
			Return(model.OwnerPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, student1.ID).
			Return(model.ReporterPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			GetAccessLevelOfUserInProject(project.ProjectID, student2.ID).
			Return(model.ReporterPermissions, nil).
			Times(1)

		gitlabRepo.EXPECT().
			ChangeUserAccessLevelInProject(project.ProjectID, student1.ID, model.DeveloperPermissions).
			Return(nil).
			Times(1)

		gitlabRepo.EXPECT().
			ChangeUserAccessLevelInProject(project.ProjectID, student2.ID, model.DeveloperPermissions).
			Return(assert.AnError).
			Times(1)

		gitlabRepo.EXPECT().
			ChangeUserAccessLevelInProject(project.ProjectID, student1.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.NoError(t, err)

		gitlabRepo.AssertExpectations(t)

		updatedAssignment, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, updatedAssignment.Closed)
		assert.NotEqual(t, newTime, *updatedAssignment.DueDate)
	})

}
