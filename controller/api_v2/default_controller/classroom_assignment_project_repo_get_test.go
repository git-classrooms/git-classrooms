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

func TestGetProjectCloneUrl(t *testing.T) {
	projectId := 1
	shh_url := "git@hs-flensburg.dev:fape2866/ci-test-project.git"
	http_url := "https://hs-flensburg.dev/fape2866/ci-test-project.git"

	expectedResponse := ProjectCloneUrlResponse{
		ProjectId:     uuid.New(),
		SshUrlToRepo:  shh_url,
		HttpUrlToRepo: http_url,
	}
	expectedResponseBody, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Could not prepare Expected Response Body")
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

	assignmentProject := &database.AssignmentProjects{
		ProjectID:     projectId,
		AssignmentID:  assignment.ID,
		TeamID:        team.ID,
		ProjectStatus: database.Accepted,
	}
	testDb.InsertAssignmentProjects(assignmentProject)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("api/v2", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetGitlabRepository(gitlabRepo)
		ctx.SetAssignmentProject(assignmentProject)
		ctx.SetGitlabProjectID(projectId)

		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Get("/api/v2/classrooms/:classroomId/projects/:projectId/repo", handler.GetProjectCloneUrls)

	t.Run("GetProjectCloneUrls - repo throws error", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectById(projectId).
			Return(nil, assert.AnError).
			Times(1)

		req := httptest.NewRequest("GET", "/api/v2/classrooms/:classroomId/projects/:projectId/repo", nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetProjectCloneUrls - /api/v2/classrooms/:classroomId/projects/:projectId/repo", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GetProjectById(projectId).
			Return(
				&model.Project{
					SSHURLToRepo:  shh_url,
					HTTPURLToRepo: http_url,
					ID:            projectId,
				},
				nil,
			).
			Times(1)

		req := httptest.NewRequest("GET", "/api/v2/classrooms/:classroomId/projects/:projectId/repo", nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(expectedResponseBody), string(body))
	})

	t.Run("GetProjectCloneUrls - assignment not accepted", func(t *testing.T) {
		assignmentProject.ProjectStatus = database.Pending
		testDb.SaveAssignmentProjects(assignmentProject)

		req := httptest.NewRequest("GET", "/api/v2/classrooms/:classroomId/projects/:projectId/repo", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
