package api

import (
	"context"
	"encoding/json"
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

func TestGetMultipleProjectCloneUrl(t *testing.T) {
	restoreDatabase(t)

	user := factory.User()
	classroom := factory.Classroom(user.ID)
	team := factory.Team(classroom.ID, make([]*database.UserClassrooms, 0))
	dueDate := time.Now().Add(1 * time.Hour)
	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	project1 := factory.AssignmentProject(assignment.ID, team.ID)
	project2 := factory.AssignmentProject(assignment.ID, team.ID)
	project3 := factory.AssignmentProject(assignment.ID, team.ID)

	project3.ProjectStatus = database.Pending

	err := query.AssignmentProjects.WithContext(context.Background()).Save(project3)
	if err != nil {
		t.Fatal("Could not save classroom!")
	}

	expectedResponse := []ProjectCloneUrlResponse{
		{
			ProjectId:     project1.ID,
			SshUrlToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project.git",
			HttpUrlToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project.git",
		},
		{
			ProjectId:     project2.ID,
			SshUrlToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project2.git",
			HttpUrlToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project2.git",
		},
	}

	expectedResponseBody, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Could not prepare Expected Response Body")
	}

	app, gitlabRepo, _ := setupApp(t, user)

	t.Run("GetProjectCloneUrls - repo throws error", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectById(project1.ID).
			Return(nil, assert.AnError).
			Times(1)

		req := httptest.NewRequest("GET", "/api/v2/classrooms/:classroomId/assignments/:assignmentId/repos", nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetProjectCloneUrls - /api/v2/classrooms/:classroomId/assignments/:assignmentId/repos", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectById(project1.ProjectID).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[0].SshUrlToRepo,
					HTTPURLToRepo: expectedResponse[0].HttpUrlToRepo,
					ID:            project1.ProjectID,
				},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetProjectById(project2.ProjectID).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[1].SshUrlToRepo,
					HTTPURLToRepo: expectedResponse[1].HttpUrlToRepo,
					ID:            project2.ProjectID,
				},
				nil,
			).
			Times(1)

		req := httptest.NewRequest("GET", "/api/v2/classrooms/:classroomId/assignments/:assignmentId/repos", nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(expectedResponseBody), string(body))
	})
}
