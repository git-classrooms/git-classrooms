package default_controller

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

func TestPutOwnedAssignments(t *testing.T) {
	testDb := db_tests.NewTestDB(t)

	user := database.User{ID: 1}
	testDb.InsertUser(&user)

	classroom := database.Classroom{
		ID:      uuid.New(),
		OwnerID: user.ID,
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
	}
	testDb.InsertTeam(&team)

	project := database.AssignmentProjects{
		AssignmentID: assignment.ID,
		TeamID:       team.ID,
		ProjectID:    1,
	}
	testDb.InsertAssignmentProjects(&project)
	assignment.Projects = append(assignment.Projects, &project)

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetOwnedClassroomAssignment(&assignment)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(user.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)
	app.Put("/api/classrooms/owned/:classroomId/assignments/:assignmentId", handler.PutOwnedAssignments)

	targetRoute := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

	t.Run("updates assignment", func(t *testing.T) {
		newTime := time.Now().Add(time.Hour * 24)
		newTime = newTime.Truncate(time.Second)
		requestBody := updateAssignmentRequest{
			Name:        "New",
			Description: "new",
			DueDate:     &newTime,
		}

		req := newPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		updatedAssignment, err := query.Assignment.WithContext(context.Background()).Where(query.Assignment.ID.Eq(assignment.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, requestBody.Name, updatedAssignment.Name)
		assert.Equal(t, requestBody.Description, updatedAssignment.Description)
		assert.Equal(t, newTime, *updatedAssignment.DueDate)
	})

	t.Run("request body is empty", func(t *testing.T) {
		requestBody := updateAssignmentRequest{}

		req := newPutJsonRequest(targetRoute, requestBody)
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

		req := newPutJsonRequest(targetRoute, requestBody)
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

		req := newPutJsonRequest(targetRoute, requestBody)
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

}
