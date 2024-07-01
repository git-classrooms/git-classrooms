package api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestGetProjectCloneUrl(t *testing.T) {
	projectId := 1
	shh_url := "git@hs-flensburg.dev:fape2866/ci-test-project.git"
	http_url := "https://hs-flensburg.dev/fape2866/ci-test-project.git"

	expectedResponse := ProjectCloneUrlResponse{
		ProjectId:     projectId,
		SshUrlToRepo:  shh_url,
		HttpUrlToRepo: http_url,
	}
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
		ctx.SetGitlabProjectID(projectId)

		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo)
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
}
