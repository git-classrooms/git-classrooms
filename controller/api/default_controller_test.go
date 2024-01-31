package api

import (
	"bytes"
	"encoding/json"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultController(t *testing.T) {
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {

		c.Locals("gitlab-repo", gitlabRepo)

		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("CreateClassroom", func(t *testing.T) {
		app.Post("/api/createClassroom", handler.CreateClassroom)

		requestBody := CreateClassroomRequest{
			Name:         "Test",
			MemberEmails: []string{"User1", "User2"},
			Description:  "test",
		}

		gitlabRepo.
			EXPECT().
			CreateGroup(
				requestBody.Name,
				model.Private,
				requestBody.Description,
				requestBody.MemberEmails,
			).
			Return(
				&model.Group{ID: 1},
				nil,
			).
			Times(1)

		for _, memberEmail := range requestBody.MemberEmails {
			gitlabRepo.
				EXPECT().
				CreateGroupInvite(1, memberEmail).
				Return(nil).
				Times(1)
		}

		req := newPostJsonRequest("/api/createClassroom", requestBody)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	})

	t.Run("CreateAssignment", func(t *testing.T) {
		app.Post("/api/createAssignment", handler.CreateAssignment)

		requestBody := CreateAssignmentRequest{
			AssigneeUserIds:   []int{12, 23},
			TemplateProjectId: 5,
		}

		template := model.Project{
			ID:   requestBody.TemplateProjectId,
			Name: "Test",
		}
		user1 := model.User{
			ID:       requestBody.AssigneeUserIds[0],
			Username: "Name1",
		}
		user2 := model.User{
			ID:       requestBody.AssigneeUserIds[1],
			Username: "Name2",
		}
		users := []model.User{user1, user2}
		name := "Test_" + user1.Username + "_" + user2.Username
		fork := model.Project{
			Name: name,
			ID:   33,
		}
		forkWithMembers := model.Project{
			Name:   fork.Name,
			ID:     fork.ID,
			Member: users,
		}

		gitlabRepo.
			EXPECT().
			GetProjectById(requestBody.TemplateProjectId).
			Return(&template, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetUserById(requestBody.AssigneeUserIds[0]).
			Return(&user1, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetUserById(requestBody.AssigneeUserIds[1]).
			Return(&user2, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ForkProject(requestBody.TemplateProjectId, name).
			Return(&fork, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			AddProjectMembers(fork.ID, users).
			Return(&forkWithMembers, nil).
			Times(1)

		req := newPostJsonRequest("/api/createAssignment", requestBody)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	})
}

func newPostJsonRequest(route string, object any) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest("POST", route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}
