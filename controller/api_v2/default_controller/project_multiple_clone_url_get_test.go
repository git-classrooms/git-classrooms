package api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestGetMultipleProjectCloneUrl(t *testing.T) {
	projectId1 := 1
	projectId2 := 2

	expectedResponse := []ProjectCloneUrlResponse{
		{
			ProjectId:     projectId1,
			SshUrlToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project.git",
			HttpUrlToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project.git",
		},
		{
			ProjectId:     projectId2,
			SshUrlToRepo:  "git@hs-flensburg.dev:fape2866/ci-test-project2.git",
			HttpUrlToRepo: "https://hs-flensburg.dev/fape2866/ci-test-project2.git",
		},
	}

	testDb := db_tests.NewTestDB(t)

	user := &database.User{
		ID:             1,
		GitlabUsername: "user1",
		GitlabEmail:    "user1",
	}
	testDb.InsertUser(user)

	classroom := &database.Classroom{
		ID:      uuid.New(),
		OwnerID: user.ID,
	}
	testDb.InsertClassroom(classroom)

	team := &database.Team{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
	}
	testDb.InsertTeam(team)

	assignment := &database.Assignment{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
	}
	testDb.InsertAssignment(assignment)

	assignmentProject1 := &database.AssignmentProjects{
		ProjectID:     projectId1,
		AssignmentID:  assignment.ID,
		TeamID:        team.ID,
		ProjectStatus: database.Accepted,
	}
	testDb.InsertAssignmentProjects(assignmentProject1)

	assignmentProject2 := &database.AssignmentProjects{
		ProjectID:     projectId2,
		AssignmentID:  assignment.ID,
		TeamID:        team.ID,
		ProjectStatus: database.Accepted,
	}
	testDb.InsertAssignmentProjects(assignmentProject2)

	assignmentProject3 := &database.AssignmentProjects{
		AssignmentID:  assignment.ID,
		TeamID:        team.ID,
		ProjectStatus: database.Pending,
	}
	testDb.InsertAssignmentProjects(assignmentProject3)

	expectedResponseBody, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Could not prepare Expected Response Body")
	}

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("api/v2", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetGitlabRepository(gitlabRepo)
		ctx.SetAssignment(assignment)

		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Get("/api/v2/classrooms/:classroomId/assignments/:assignmentId/repos", handler.GetMultipleProjectCloneUrls)

	t.Run("GetProjectCloneUrls - repo throws error", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectById(projectId1).
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
			GetProjectById(projectId1).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[0].SshUrlToRepo,
					HTTPURLToRepo: expectedResponse[0].HttpUrlToRepo,
					ID:            projectId1,
				},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetProjectById(projectId2).
			Return(
				&model.Project{
					SSHURLToRepo:  expectedResponse[1].SshUrlToRepo,
					HTTPURLToRepo: expectedResponse[1].HttpUrlToRepo,
					ID:            projectId2,
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
