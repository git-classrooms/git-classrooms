package default_controller

import (
	"context"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassrooms(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(user)

	classrooms := []*database.Classroom{
		{
			Name:               "Classroom One",
			OwnerID:            1,
			Description:        "Description One",
			GroupID:            10,
			GroupAccessTokenID: 30,
			GroupAccessToken:   "token30",
		},
		{
			Name:               "Classroom Two",
			OwnerID:            1,
			Description:        "Description Two",
			GroupID:            20,
			GroupAccessTokenID: 40,
			GroupAccessToken:   "token40",
		},
	}

	for _, classroom := range classrooms {
		testDB.InsertClassroom(classroom)
	}

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassrooms", func(t *testing.T) {
		app.Get("/api/classrooms/owned", handler.GetOwnedClassrooms)
		req := httptest.NewRequest("GET", "/api/classrooms/owned", nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		classrooms, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(1)).Find()
		assert.NoError(t, err)
		assert.Len(t, classrooms, len(classrooms))
		for i, classroom := range classrooms {
			assert.Equal(t, classrooms[i].Name, classroom.Name)
			assert.Equal(t, classrooms[i].Description, classroom.Description)
			assert.Equal(t, classrooms[i].GroupID, classroom.GroupID)
			assert.Equal(t, classrooms[i].GroupAccessTokenID, classroom.GroupAccessTokenID)
			assert.Equal(t, classrooms[i].GroupAccessToken, classroom.GroupAccessToken)
		}
	})
}

