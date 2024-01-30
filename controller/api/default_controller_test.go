package api

import (
	"bytes"
	"de.hs-flensburg.gitlab/gitlab-classroom/model"
	gitlabRepoMock "de.hs-flensburg.gitlab/gitlab-classroom/repository/gitlab/_mock"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiHandler(t *testing.T) {
	// ctrl := gomock.NewController(t)
	repo := gitlabRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {

		c.Locals("gitlab-repo", repo)

		return c.Next()
	})

	handler := NewApiController()

	t.Run("CreateClassroom", func(t *testing.T) {
		app.Post("/api/createClassroom", handler.CreateClassroom)

		requestBody := CreateClassroomRequest{
			Name:         "Test",
			MemberEmails: []string{"User1", "User2"},
			Description:  "test",
		}

		repo.
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
			repo.
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

		repo.
			EXPECT().
			GetProjectById(requestBody.TemplateProjectId).
			Return(&template, nil).
			Times(1)

		repo.
			EXPECT().
			GetUserById(requestBody.AssigneeUserIds[0]).
			Return(&user1, nil).
			Times(1)

		repo.
			EXPECT().
			GetUserById(requestBody.AssigneeUserIds[1]).
			Return(&user2, nil).
			Times(1)

		repo.
			EXPECT().
			ForkProject(requestBody.TemplateProjectId, name).
			Return(&fork, nil).
			Times(1)

		repo.
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
