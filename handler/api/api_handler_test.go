package apiHandler

import (
	mock_repository "backend/api/repository/_mocks"
	"backend/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApiHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockRepository(ctrl)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {

		c.Locals("gitlab-repo", repo)

		return c.Next()
	})

	handler := NewFiberApiHandler()

	t.Run("CreateClassroom", func(t *testing.T) {
		app.Post("/api/createClassroom", handler.CreateClassroom)

		requestBody := ClassroomRequest{
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

		resp, err := app.Test(req, 1)

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
