package default_controller

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	databaseConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/database"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateClassroom(t *testing.T) {
	// --------------- DB SETUP -----------------
	pq, err := tests.StartPostgres()
	t.Cleanup(func() {
		pq.Terminate(context.Background())
	})
	port, err := pq.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("could not get database container port: %s", err.Error())
	}
	dbConfig := databaseConfig.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	}
	db, err := gorm.Open(postgres.Open(dbConfig.Dsn()), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to database: %s", err.Error())
	}
	err = utils.MigrateDatabase(db)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}

	query.SetDefault(db)

	insertTestData(t)
	// ------------ END OF DB SETUP -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		fiberContext.SetGitlabRepository(c, gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("CreateClassroom", func(t *testing.T) {
		app.Post("/api/classrooms", handler.CreateClassroom)

		requestBody := CreateClassroomRequest{
			Name:         "Test",
			MemberEmails: []string{},
			Description:  "test",
		}

		gitlabRepo.
			EXPECT().
			CreateGroup(
				requestBody.Name,
				model.Private,
				requestBody.Description,
			).
			Return(
				&model.Group{ID: 1},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			CreateGroupAccessToken(
				1,
				"Gitlab Classrooms",
				model.OwnerPermissions,
				mock.AnythingOfType("time.Time"),
				"api",
			).
			Return(
				&model.GroupAccessToken{ID: 20, Token: "token"},
				nil,
			).
			Times(1)

		req := newPostJsonRequest("/api/classrooms", requestBody)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		assert.NoError(t, err)

		classRoom, err := query.Classroom.Where(query.Classroom.OwnerID.Eq(1)).First()
		assert.NoError(t, err)
		assert.Equal(t, "Test", classRoom.Name)
		assert.Equal(t, "test", classRoom.Description)
		assert.Equal(t, 1, classRoom.GroupID)
		assert.Equal(t, 20, classRoom.GroupAccessTokenID)
		assert.Equal(t, "token", classRoom.GroupAccessToken)

		assert.Equal(t, fmt.Sprintf("/api/v1/classrooms/%s", classRoom.ID.String()), resp.Header.Get("Location"))
	})
}

func insertTestData(t *testing.T) {
	user := &database.User{ID: 1}
	err := query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not save user: %s", err.Error())
	}
}
