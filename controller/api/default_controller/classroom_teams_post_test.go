package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestCreateTeam(t *testing.T) {
	// setup database
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	dueDate := time.Now().Add(1 * time.Hour)

	factory.Assignment(classroom.ID, &dueDate, false)
	teamGroupId := gofakeit.Int()

	app, gitlabRepo, _ := setupApp(t, owner)

	t.Run("TestCreateTeam", func(t *testing.T) {
		requestBody := createTeamRequest{
			Name: gofakeit.Name(),
		}

		gitlabRepo.
			EXPECT().
			GroupAccessLogin(classroom.GroupAccessToken).
			Return(nil).
			Times(1)

		teamDescription := fmt.Sprintf("Team %s of classroom %s", requestBody.Name, classroom.Name)
		gitlabRepo.
			EXPECT().
			CreateSubGroup(
				requestBody.Name,
				requestBody.Name,
				classroom.GroupID,
				model.Private,
				teamDescription,
			).
			Return(&model.Group{Name: requestBody.Name, ID: teamGroupId}, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeGroupDescription(teamGroupId, mock.Anything).
			Return(nil, nil).
			Times(1)

		route := fmt.Sprintf("/api/v1/classrooms/%s/teams", classroom.ID)

		req := newPostJsonRequest(route, requestBody)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assert.NoError(t, err)

		queryTeam := query.Team

		team, err := queryTeam.
			WithContext(context.Background()).
			First()
		assert.NoError(t, err)

		assert.Equal(t, team.Name, requestBody.Name)
	})
}
