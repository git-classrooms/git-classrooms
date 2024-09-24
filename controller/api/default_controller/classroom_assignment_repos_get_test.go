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

func TestGetMultipleProjectCloneUrl(t *testing.T) {
	restoreDatabase(t)

	user := factory.User()
	classroom := factory.Classroom(user.ID)
	factory.UserClassroom(user.ID, classroom.ID, database.Owner)

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

	expectedResponse := []ProjectCloneURLResponse{
		{
			ProjectID:     project1.ID,
			SSHURLToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project.git",
			HTTPURLToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project.git",
		},
		{
			ProjectID:     project2.ID,
			SSHURLToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project2.git",
			HTTPURLToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project2.git",
		},
	}

	expectedResponseBody, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Could not prepare Expected Response Body")
	}

	app, gitlabRepo, _ := setupApp(t, user)
	targetRoute := fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s/repos", classroom.ID, assignment.ID)

	t.Run("GetProjectCloneUrls - repo throws error", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectByID(project1.ProjectID).
			Return(nil, assert.AnError).
			Times(1)

		req := httptest.NewRequest("GET", targetRoute, nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetProjectCloneUrls", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectByID(project1.ProjectID).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[0].SSHURLToRepo,
					HTTPURLToRepo: expectedResponse[0].HTTPURLToRepo,
					ID:            project1.ProjectID,
				},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetProjectByID(project2.ProjectID).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[1].SSHURLToRepo,
					HTTPURLToRepo: expectedResponse[1].HTTPURLToRepo,
					ID:            project2.ProjectID,
				},
				nil,
			).
			Times(1)

		req := httptest.NewRequest("GET", targetRoute, nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(expectedResponseBody), string(body))
	})
}
