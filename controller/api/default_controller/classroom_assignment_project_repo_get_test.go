package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestGetProjectCloneUrl(t *testing.T) {
	restoreDatabase(t)

	shhURL := "git@hs-flensburg.dev:fape2866/ci-test-project.git"
	httpURL := "https://hs-flensburg.dev/fape2866/ci-test-project.git"

	user := factory.User()
	classroom := factory.Classroom(user.ID)
	factory.UserClassroom(user.ID, classroom.ID, database.Owner)

	team := factory.Team(classroom.ID, make([]*database.UserClassrooms, 0))

	dueDate := time.Now().Add(1 * time.Hour)

	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	assignmentProject := factory.AssignmentProject(assignment.ID, team.ID)

	expectedResponse := ProjectCloneURLResponse{
		ProjectID:     assignmentProject.ID,
		SSHURLToRepo:  shhURL,
		HTTPURLToRepo: httpURL,
	}
	expectedResponseBody, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Could not prepare Expected Response Body")
	}
	app, gitlabRepo, _ := setupApp(t, user)

	targetRoute := fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s/projects/%s/repo", classroom.ID, assignment.ID, assignmentProject.ID)

	t.Run("GetProjectCloneUrls - repo throws error", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectByID(assignmentProject.ProjectID).
			Return(nil, assert.AnError).
			Times(1)

		req := httptest.NewRequest("GET", targetRoute, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetProjectCloneUrls - /api/v1/classrooms/:classroomId/projects/:projectId/repo", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectByID(assignmentProject.ProjectID).
			Return(
				&model.Project{
					SSHURLToRepo:  shhURL,
					HTTPURLToRepo: httpURL,
					ID:            assignmentProject.ProjectID,
				},
				nil,
			).
			Times(1)

		req := httptest.NewRequest("GET", targetRoute, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(expectedResponseBody), string(body))
	})

	t.Run("GetProjectCloneUrls - assignment not accepted", func(t *testing.T) {
		assignmentProject.ProjectStatus = database.Pending
		err := query.AssignmentProjects.WithContext(context.Background()).Save(assignmentProject)
		if err != nil {
			t.Fatalf("could not update assignment project: %s", err.Error())
		}

		req := httptest.NewRequest("GET", targetRoute, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
