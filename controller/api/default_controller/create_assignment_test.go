package default_controller

import (
	"testing"
)

func TestCreateAssignment(t *testing.T) {
	/*
		gitlabRepo := gitlabRepoMock.NewMockRepository(t)
		mailRepo := mailRepoMock.NewMockRepository(t)

		app := fiber.New()
		app.Use("/api", func(c *fiber.Ctx) error {

			c.Locals("gitlab-repo", gitlabRepo)

			return c.Next()
		})

		handler := NewApiController(mailRepo)

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
		})*/
}
