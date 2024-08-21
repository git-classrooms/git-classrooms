package api

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestPutOwnedAssignments(t *testing.T) {
	testDb := db_tests.NewTestDB(t)

	owner := database.User{
		ID:             1,
		GitlabUsername: "owner",
		GitlabEmail:    "owner",
	}
	testDb.InsertUser(&owner)

	student1 := database.User{
		ID:             2,
		GitlabUsername: "student1",
		GitlabEmail:    "student1",
	}
	testDb.InsertUser(&student1)

	student2 := database.User{
		ID:             3,
		GitlabUsername: "student2",
		GitlabEmail:    "student2",
	}
	testDb.InsertUser(&student2)

	classroom := database.Classroom{
		ID:      uuid.New(),
		OwnerID: owner.ID,
	}
	testDb.InsertClassroom(&classroom)
	oldTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	oldTime = oldTime.Truncate(time.Second)
	assignment := database.Assignment{
		ID:          uuid.New(),
		Name:        "Test",
		Description: "test",
		DueDate:     &oldTime,
		ClassroomID: classroom.ID,
	}
	testDb.InsertAssignment(&assignment)

	team := database.Team{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
		GroupID:     10,
		Member: []*database.UserClassrooms{
			{
				UserID:      owner.ID,
				ClassroomID: classroom.ID,
				Role:        database.Owner,
			},
			{
				UserID:      student1.ID,
				ClassroomID: classroom.ID,
				Role:        database.Student,
			},
			{
				UserID:      student2.ID,
				ClassroomID: classroom.ID,
				Role:        database.Student,
			},
		},
	}
	testDb.InsertTeam(&team)

	project := database.AssignmentProjects{
		AssignmentID: assignment.ID,
		TeamID:       team.ID,
		ProjectID:    1,
	}
	testDb.InsertAssignmentProjects(&project)
	assignment.Projects = append(assignment.Projects, &project)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetAssignment(&assignment)
		ctx.SetGitlabRepository(gitlabRepo)

		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Put("/api/classrooms/:classroomId/assignments/:assignmentId", handler.UpdateAssignment)

	targetRoute := fmt.Sprintf("/api/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

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

	t.Run("request body is empty", func(t *testing.T) {
		requestBody := updateAssignmentRequest{}

		req := db_tests.NewPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		assert.Equal(t, "Request can not be empty, requires name, description or dueDate", bodyString)
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
		testDb.SaveAssignmentProjects(&project)

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
		testDb.SaveAssignmentProjects(&project)

		assignment.Closed = true
		testDb.SaveAssignment(&assignment)

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
		testDb.SaveAssignmentProjects(&project)

		assignment.Closed = true
		testDb.SaveAssignment(&assignment)

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
