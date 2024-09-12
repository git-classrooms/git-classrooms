package default_controller

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetMe(t *testing.T) {
	sqlMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An unknown error accured, %v", err)
	}
	_ = mock
	defer sqlMock.Close()

	postgresDB := postgres.New(postgres.Config{Conn: sqlMock})

	db, err := gorm.Open(postgresDB)
	if err != nil {
		t.Fatalf("An unknown error accured, %v", err)
	}

	db = db.Debug()

	query.SetDefault(db)

	session.InitSessionStore(nil, &url.URL{Scheme: "http"})
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("Get Me", func(t *testing.T) {
		app.Get("/api/v1/me", handler.GetMe)
		route := "/api/v1/me"

		shouldUser := database.User{
			ID:          1,
			Name:        "Toni Tester",
			GitlabEmail: "toni@tester.com",
		}

		avatarUrl := "https://gitlab.com/avatar.png"
		shouldAvatar := database.UserAvatar{
			UserID:            1,
			AvatarURL:         &avatarUrl,
			FallbackAvatarURL: &avatarUrl,
		}

		userSqlRows := GetUserSqlMockRows(shouldUser)
		userAvatarSqlRows := GetUserAvatarSqlMockRows(shouldAvatar)

		mock.ExpectQuery("^SELECT (.+) FROM \"users\" (.+)$").WillReturnRows(userSqlRows)
		mock.ExpectQuery("^SELECT (.+) FROM \"user_avatars\" (.+)$").WillReturnRows(userAvatarSqlRows)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func GetUserSqlMockRows(users ...database.User) *sqlmock.Rows {
	sqlRows := sqlmock.NewRows([]string{"id", "gitlab_email", "name"})

	for _, u := range users {
		sqlRows.AddRow(u.ID, u.Name, u.GitlabEmail)
	}
	return sqlRows
}

func GetUserAvatarSqlMockRows(avatars ...database.UserAvatar) *sqlmock.Rows {
	sqlRows := sqlmock.NewRows([]string{"user_id", "avatar_url", "fallback_avatar_url"})

	for _, a := range avatars {
		sqlRows.AddRow(a.UserID, a.AvatarURL, a.FallbackAvatarURL)
	}
	return sqlRows
}
